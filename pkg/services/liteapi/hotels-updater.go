package liteapi

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"sync"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/go-co-op/gocron/v2"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lockerRepoInstance repo.LockerRepo

type chData struct {
	ctx   context.Context
	place models.Place
}

type locker struct{}

func (l *locker) Lock(ctx context.Context, key string) (gocron.Lock, error) {

	if lockerRepoInstance == nil {
		return nil, errors.New("locker repo instance is not set")
	}

	//todo : acquire the lock
	if err := lockerRepoInstance.EnsureIndexes(ctx); err != nil {
		fmt.Println("error create index")
		return nil, err
	}

	// Allow acquiring locks that are either unlocked or expired
	filter := bson.M{
		"key": key,
		"$or": []bson.M{
			{"isLocked": false},
			{"expiresAt": bson.M{"$lt": time.Now()}},
		},
	}

	expiresAt := time.Now().Add(2 * time.Hour)
	update := bson.M{
		"$set": bson.M{
			"isLocked":  true,
			"expiresAt": expiresAt,
		},
	}

	upsert := true
	after := options.After
	opts := []*options.FindOneAndUpdateOptions{
		{
			Upsert:         &upsert,
			ReturnDocument: &after,
		},
	}

	doc, err := lockerRepoInstance.Patch(ctx, filter, update, opts...)
	if err != nil {
		// if mongo.IsDuplicateKeyError(err) {
		// 	return nil, nil
		// }
		fmt.Println("faild to acquire the lock", err)
		return nil, err
	}

	fmt.Println("lock acquired successfully for", key)

	return &acquiredLock{
		lock: doc,
	}, nil
}

type acquiredLock struct {
	lock *models.Locker
}

func (al *acquiredLock) Unlock(ctx context.Context) error {
	if lockerRepoInstance == nil {
		return errors.New("locker repo instance is not set to unlock the lock")
	}
	//todo : release the lock
	filter := bson.M{
		"_id": al.lock.Id,
	}

	update := bson.M{
		"$set": bson.M{"isLocked": false},
	}

	_, err := lockerRepoInstance.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

type HotelsUpdater struct {
	lockerRepo  repo.LockerRepo
	placeRepo   repo.PlaceRepo
	dataSvcs    DataSvcs
	numOfWorker int
	dbs         []string
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	workChan    chan chData
}

func NewHotelsUpdater(i *do.Injector, dbs []string) *HotelsUpdater {
	ctx, cancel := context.WithCancel(context.Background())
	return &HotelsUpdater{
		lockerRepo:  do.MustInvoke[repo.LockerRepo](i),
		placeRepo:   do.MustInvoke[repo.PlaceRepo](i),
		dataSvcs:    do.MustInvoke[DataSvcs](i),
		numOfWorker: 20,
		dbs:         dbs,
		ctx:         ctx,
		cancel:      cancel,
		workChan:    make(chan chData),
	}
}

func (u *HotelsUpdater) worker(c chan chData, i int) {
	defer func() {
		u.wg.Done()
	}()

	for {
		select {
		case <-u.ctx.Done():
			fmt.Println("worker", i, "stopping due to context cancellation")
			return
		case data, ok := <-c:
			if !ok {
				fmt.Println("worker", i, "channel closed, stopping")
				return
			}
			place := data.place
			ctx := data.ctx
			cfg, _ := util.GetReqAppCfg(ctx)

			fmt.Printf("executing place: %s , db %s \n", place.PlaceId, cfg.Db)
			query := map[string]string{
				"language":      place.Language,
				"placeId":       place.PlaceId,
				"lastUpdatedAt": time.Now().Add(-2 * 24 * time.Hour).Format(time.RFC3339),
			}

			fmt.Println(query)

			u.dataSvcs.GetHotelsByPlaceId(ctx, query)

		}
	}

}

func (u *HotelsUpdater) updateHotels(c chan chData) error {
	for _, db := range u.dbs {

		cfg := &types.AppCfg{
			Db: db,
			Hp: &common.HeaderParams{
				Client: db,
			},
		}

		ctx := context.WithValue(context.Background(), util.ReqAppCfg, cfg)

		pipeline := []bson.M{
			{"$match": bson.M{}},
		}

		var places []models.Place
		if err := u.placeRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &places)
		}); err != nil {
			fmt.Println("error aggregate places", err)
			return err
		}

		for _, place := range places {
			select {
			case <-u.ctx.Done():
				fmt.Println("updateHotels stopping due to context cancellation")
				return nil
			case c <- chData{
				ctx:   ctx,
				place: place,
			}:
			}
		}
	}

	return nil

}

var scheduler gocron.Scheduler

func (u *HotelsUpdater) Run() error {

	var err error

	lockerRepoInstance = u.lockerRepo

	u.workChan = make(chan chData)

	for i := 0; i < u.numOfWorker; i++ {
		u.wg.Add(1)
		go u.worker(u.workChan, i)
	}

	locker := &locker{}
	// create a scheduler
	scheduler, err = gocron.NewScheduler(
		gocron.WithDistributedLocker(locker),
	)
	if err != nil {
		return err
	}

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			48*time.Hour,
		),
		gocron.NewTask(
			u.updateHotels,
			u.workChan,
		),
	)
	if err != nil {
		return err
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	fmt.Println("starting scheduler")
	scheduler.Start()

	<-u.ctx.Done()

	fmt.Println("updater stopped")

	return nil
}

func (u *HotelsUpdater) Stop() {
	fmt.Println("stopping scheduler...")

	// Cancel context first to stop new job schedules
	u.cancel()

	// Wait a bit for current jobs to finish gracefully
	time.Sleep(2 * time.Second)

	// Stop the scheduler to prevent new job executions
	if scheduler != nil {
		if err := scheduler.Shutdown(); err != nil {
			fmt.Println("error shutting down scheduler", err)
		} else {
			fmt.Println("scheduler shut down successfully")
		}
	}

	// Close the channel to signal workers
	if u.workChan != nil {
		close(u.workChan)
	}

	// Wait for workers to finish
	fmt.Println("waiting for workers to finish...")
	u.wg.Wait()
	fmt.Println("all workers stopped")

	// Clean up any remaining locks
	u.cleanupLocks()
}

func (u *HotelsUpdater) cleanupLocks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("cleaning up remaining locks...")

	// Release all locks that might be stuck
	filter := bson.M{"isLocked": true}
	update := bson.M{"$set": bson.M{"isLocked": false}}

	// Use BulkWrite to update all locked documents
	var updateModels []mongo.WriteModel

	// First, get all locked documents
	pipeline := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"_id": 1}},
	}

	var lockedDocs []struct {
		ID interface{} `bson:"_id"`
	}

	err := u.lockerRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &lockedDocs)
	})

	if err != nil {
		fmt.Println("error finding locked documents:", err)
		return
	}

	// Create update operations for each locked document
	for _, doc := range lockedDocs {
		updateModel := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": doc.ID}).
			SetUpdate(update)
		updateModels = append(updateModels, updateModel)
	}

	if len(updateModels) == 0 {
		fmt.Println("no locks to clean up")
		return
	}

	// Execute bulk write
	_, err = u.lockerRepo.BulkWrite(ctx, updateModels)
	if err != nil {
		fmt.Println("error cleaning up locks:", err)
	} else {
		fmt.Printf("successfully cleaned up %d lock(s)\n", len(updateModels))
	}
}
