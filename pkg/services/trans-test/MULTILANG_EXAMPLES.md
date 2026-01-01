# Multi-Language TransTest Solution

This document provides examples and usage patterns for the multi-language TransTest implementation.

## Features

- **Efficient Storage**: Translations are stored as embedded maps within the same document
- **Automatic Fallback**: Supports fallback to default languages (en, ar) when requested language is not available
- **Reusable Translation Type**: The `Translation` type can be used across other models
- **Language-Specific Queries**: Get localized responses without client-side processing
- **Granular Updates**: Update individual translations without affecting others

## JSON Examples

### 1. Creating Multi-Language Data (TransTestDto)

```json
{
  "name": {
    "en": "Tourism Expert John Doe",
    "ar": "خبير السياحة جون دو",
    "fr": "Expert en tourisme John Doe",
    "es": "Experto en turismo John Doe"
  },
  "description": {
    "en": "Senior Tourism Consultant with 10+ years of experience in Middle East tourism",
    "ar": "مستشار سياحي كبير مع أكثر من 10 سنوات من الخبرة في سياحة الشرق الأوسط",
    "fr": "Consultant senior en tourisme avec plus de 10 ans d'expérience dans le tourisme au Moyen-Orient",
    "es": "Consultor senior en turismo con más de 10 años de experiencia en turismo de Oriente Medio"
  },
  "about": {
    "bio": {
      "en": "Passionate about creating unforgettable travel experiences across the Middle East. Specialized in luxury tourism, cultural tours, and adventure travel.",
      "ar": "شغوف بخلق تجارب سفر لا تُنسى عبر الشرق الأوسط. متخصص في السياحة الفاخرة والجولات الثقافية وسياحة المغامرات.",
      "fr": "Passionné par la création d'expériences de voyage inoubliables au Moyen-Orient. Spécialisé dans le tourisme de luxe, les circuits culturels et les voyages d'aventure.",
      "es": "Apasionado por crear experiencias de viaje inolvidables en Oriente Medio. Especializado en turismo de lujo, tours culturales y viajes de aventura."
    },
    "testimonials": [
      {
        "name": {
          "en": "Sarah Johnson",
          "ar": "سارة جونسون",
          "fr": "Sarah Johnson",
          "es": "Sarah Johnson"
        },
        "testimonial": {
          "en": "John provided exceptional service and made our trip to Jordan absolutely perfect!",
          "ar": "قدم جون خدمة استثنائية وجعل رحلتنا إلى الأردن مثالية تماماً!",
          "fr": "John a fourni un service exceptionnel et a rendu notre voyage en Jordanie absolument parfait!",
          "es": "¡John brindó un servicio excepcional e hizo que nuestro viaje a Jordania fuera absolutamente perfecto!"
        }
      },
      {
        "name": {
          "en": "Ahmed Al-Rashid",
          "ar": "أحمد الراشد",
          "fr": "Ahmed Al-Rashid",
          "es": "Ahmed Al-Rashid"
        },
        "testimonial": {
          "en": "Professional, knowledgeable, and always available. Highly recommended!",
          "ar": "محترف ومطلع ومتاح دائماً. أنصح به بشدة!",
          "fr": "Professionnel, compétent et toujours disponible. Hautement recommandé!",
          "es": "Profesional, conocedor y siempre disponible. ¡Altamente recomendado!"
        }
      }
    ],
    "project": {
      "project_name": {
        "en": "Middle East Cultural Heritage Tours",
        "ar": "جولات التراث الثقافي للشرق الأوسط",
        "fr": "Circuits du patrimoine culturel du Moyen-Orient",
        "es": "Tours del patrimonio cultural de Oriente Medio"
      },
      "project_description": {
        "en": "Comprehensive tourism platform connecting travelers with authentic Middle Eastern experiences",
        "ar": "منصة سياحية شاملة تربط المسافرين بتجارب أصيلة من الشرق الأوسط",
        "fr": "Plateforme touristique complète reliant les voyageurs aux expériences authentiques du Moyen-Orient",
        "es": "Plataforma turística integral que conecta a los viajeros con experiencias auténticas de Oriente Medio"
      },
      "details": {
        "description": {
          "en": "Built with cutting-edge technology to provide seamless booking, real-time updates, and personalized recommendations",
          "ar": "مبني بأحدث التقنيات لتوفير حجز سلس وتحديثات فورية وتوصيات مخصصة",
          "fr": "Construit avec une technologie de pointe pour fournir une réservation transparente, des mises à jour en temps réel et des recommandations personnalisées",
          "es": "Construido con tecnología de vanguardia para proporcionar reservas sin problemas, actualizaciones en tiempo real y recomendaciones personalizadas"
        },
        "location": {
          "en": "Jordan, UAE, Saudi Arabia",
          "ar": "الأردن، الإمارات العربية المتحدة، المملكة العربية السعودية",
          "fr": "Jordanie, Émirats arabes unis, Arabie saoudite",
          "es": "Jordania, Emiratos Árabes Unidos, Arabia Saudita"
        }
      }
    }
  },
  "notes": [
    {
      "en": "Certified by Jordan Tourism Board",
      "ar": "معتمد من هيئة تنشيط السياحة الأردنية",
      "fr": "Certifié par l'Office du tourisme de Jordanie",
      "es": "Certificado por la Junta de Turismo de Jordania"
    },
    {
      "en": "Fluent in Arabic, English, and French",
      "ar": "يتقن العربية والإنجليزية والفرنسية",
      "fr": "Maîtrise l'arabe, l'anglais et le français",
      "es": "Domina el árabe, inglés y francés"
    },
    {
      "en": "Specializes in archaeological site tours",
      "ar": "متخصص في جولات المواقع الأثرية",
      "fr": "Spécialisé dans les visites de sites archéologiques",
      "es": "Se especializa en tours de sitios arqueológicos"
    }
  ]
}
```

### 2. Localized Response (English)

When calling `GetOneLocalized(ctx, id, "en")`:

```json
{
  "_id": "6507f1f77bcf86cd99439011",
  "name": "Tourism Expert John Doe",
  "description": "Senior Tourism Consultant with 10+ years of experience in Middle East tourism",
  "about": {
    "bio": "Passionate about creating unforgettable travel experiences across the Middle East. Specialized in luxury tourism, cultural tours, and adventure travel.",
    "testimonials": [
      {
        "name": "Sarah Johnson",
        "testimonial": "John provided exceptional service and made our trip to Jordan absolutely perfect!"
      },
      {
        "name": "Ahmed Al-Rashid",
        "testimonial": "Professional, knowledgeable, and always available. Highly recommended!"
      }
    ],
    "project": {
      "project_name": "Middle East Cultural Heritage Tours",
      "project_description": "Comprehensive tourism platform connecting travelers with authentic Middle Eastern experiences",
      "details": {
        "description": "Built with cutting-edge technology to provide seamless booking, real-time updates, and personalized recommendations",
        "location": "Jordan, UAE, Saudi Arabia"
      }
    }
  },
  "notes": [
    "Certified by Jordan Tourism Board",
    "Fluent in Arabic, English, and French",
    "Specializes in archaeological site tours"
  ],
  "status": "active",
  "trash": false,
  "created_at": "2023-09-18T10:30:00Z",
  "updated_at": "2023-09-18T10:30:00Z"
}
```

### 3. Localized Response (Arabic)

When calling `GetOneLocalized(ctx, id, "ar")`:

```json
{
  "_id": "6507f1f77bcf86cd99439011",
  "name": "خبير السياحة جون دو",
  "description": "مستشار سياحي كبير مع أكثر من 10 سنوات من الخبرة في سياحة الشرق الأوسط",
  "about": {
    "bio": "شغوف بخلق تجارب سفر لا تُنسى عبر الشرق الأوسط. متخصص في السياحة الفاخرة والجولات الثقافية وسياحة المغامرات.",
    "testimonials": [
      {
        "name": "سارة جونسون",
        "testimonial": "قدم جون خدمة استثنائية وجعل رحلتنا إلى الأردن مثالية تماماً!"
      },
      {
        "name": "أحمد الراشد",
        "testimonial": "محترف ومطلع ومتاح دائماً. أنصح به بشدة!"
      }
    ],
    "project": {
      "project_name": "جولات التراث الثقافي للشرق الأوسط",
      "project_description": "منصة سياحية شاملة تربط المسافرين بتجارب أصيلة من الشرق الأوسط",
      "details": {
        "description": "مبني بأحدث التقنيات لتوفير حجز سلس وتحديثات فورية وتوصيات مخصصة",
        "location": "الأردن، الإمارات العربية المتحدة، المملكة العربية السعودية"
      }
    }
  },
  "notes": [
    "معتمد من هيئة تنشيط السياحة الأردنية",
    "يتقن العربية والإنجليزية والفرنسية",
    "متخصص في جولات المواقع الأثرية"
  ],
  "status": "active",
  "trash": false,
  "created_at": "2023-09-18T10:30:00Z",
  "updated_at": "2023-09-18T10:30:00Z"
}
```

## Usage Examples

### Service Method Usage

```go
// Get localized single item
localizedItem, err := service.GetOneLocalized(ctx, "6507f1f77bcf86cd99439011", "ar", "en")

// Get localized paginated results
localizedResults, err := service.GetLocalized(ctx, 0, 10, query, "fr", "en")

// Update a specific translation
err := service.UpdateTranslation(ctx, "6507f1f77bcf86cd99439011", "name", "es", "Experto en turismo Juan Pérez")

// Add multiple translations at once
translations := map[string]string{
    "en": "New English value",
    "ar": "قيمة جديدة بالعربية",
    "fr": "Nouvelle valeur en français",
}
err := service.AddTranslation(ctx, "6507f1f77bcf86cd99439011", "description", translations)
```

### Building Translation Data

```go
// Create translation map
nameTranslations := BuildTranslation(map[string]string{
    "en": "John Doe",
    "ar": "جون دو",
    "fr": "Jean Dupont",
})

// Create translated notes
notesData := []map[string]string{
    {
        "en": "Expert in archaeology",
        "ar": "خبير في علم الآثار",
    },
    {
        "en": "Licensed tour guide",
        "ar": "دليل سياحي مرخص",
    },
}
translatedNotes := BuildTranslatedNotes(notesData)
```

## Key Benefits

1. **Single Database Query**: All translations stored in one document
2. **Automatic Fallback**: Falls back to English or Arabic if requested language not available
3. **Type Safety**: Strong typing with Go structs
4. **Easy Updates**: Update individual translations without affecting others
5. **Flexible Querying**: Get either full translation maps or localized responses
6. **Reusable**: Translation type can be used across different models

## Migration Strategy

When migrating existing data:

1. Convert existing string fields to Translation maps
2. Use English as the default language for existing data
3. Gradually add translations for other languages
4. Update frontend to use localized endpoints

```go
// Example migration helper
func migrateExistingData(existingName string) models.Translation {
    return models.Translation{
        "en": existingName, // Use existing data as English
    }
}
``` 