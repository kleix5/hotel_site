package service

import "time"

// SiteContent stores top-level homepage content and contact metadata.
type SiteContent struct {
	// ID is the database primary key.
	ID int64 `json:"id"`
	// BrandName is the visual hotel/project name in the header area.
	BrandName string `json:"brandName"`
	// Tagline is a short branding phrase shown in multiple hero/header spots.
	Tagline string `json:"tagline"`
	// HeroTitle is the main first-screen heading.
	HeroTitle string `json:"heroTitle"`
	// HeroSubtitle is the supporting copy below the hero title.
	HeroSubtitle string `json:"heroSubtitle"`
	// HeroImagePath points to optional uploaded hero image; empty means fallback art.
	HeroImagePath string `json:"heroImagePath"`
	// SeasonNote highlights operational season constraints or seasonal offer context.
	SeasonNote string `json:"seasonNote"`
	// BookingURL is CTA target for the booking action.
	BookingURL string `json:"bookingUrl"`
	// Address is the public location line shown in contacts.
	Address string `json:"address"`
	// Phone is the public booking/contact phone number.
	Phone string `json:"phone"`
	// Email is the public contact email.
	Email string `json:"email"`
	// CheckIn is the advertised check-in time.
	CheckIn string `json:"checkIn"`
	// CheckOut is the advertised check-out time.
	CheckOut string `json:"checkOut"`
	// AboutTitle is the heading of the "about concept" block.
	AboutTitle string `json:"aboutTitle"`
	// AboutBody is the long-form descriptive paragraph about the property.
	AboutBody string `json:"aboutBody"`
	// CreatedAt is insertion timestamp managed by MySQL.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is last modification timestamp managed by MySQL.
	UpdatedAt time.Time `json:"updatedAt"`
}

// Room describes one room-card entry rendered in the rooms section.
type Room struct {
	// ID is the database primary key.
	ID int64 `json:"id"`
	// Name is the room category title shown on card.
	Name string `json:"name"`
	// Summary is descriptive copy for room category.
	Summary string `json:"summary"`
	// Capacity is maximum guest count.
	Capacity int `json:"capacity"`
	// SizeSQM is approximate area in square meters.
	SizeSQM int `json:"sizeSqm"`
	// RoomCount indicates how many rooms belong to this category.
	RoomCount int `json:"roomCount"`
	// FeatureNote is a short highlight/badge shown above room title.
	FeatureNote string `json:"featureNote"`
	// ImagePath points to optional uploaded room image.
	ImagePath string `json:"imagePath"`
}

// GalleryItem describes a visual card in the gallery section.
type GalleryItem struct {
	// ID is the database primary key.
	ID int64 `json:"id"`
	// Title is the gallery card heading.
	Title string `json:"title"`
	// Subtitle is supporting text under the gallery title.
	Subtitle string `json:"subtitle"`
	// Tone selects fallback gradient class when image is absent.
	Tone string `json:"tone"`
	// ImagePath points to optional uploaded gallery image.
	ImagePath string `json:"imagePath"`
}

// Snapshot stores scraped source material metadata for reference display.
type Snapshot struct {
	// ID is the database primary key.
	ID int64 `json:"id"`
	// SourceURL is original scraped page URL.
	SourceURL string `json:"sourceUrl"`
	// PageTitle is title extracted from source page.
	PageTitle string `json:"pageTitle"`
	// Description is meta description or summary extracted during scrape.
	Description string `json:"description"`
	// Markdown is normalized markdown representation of source content.
	Markdown string `json:"markdown"`
	// HTML is captured raw/processed HTML snapshot.
	HTML string `json:"html"`
	// CreatedAt is snapshot creation timestamp.
	CreatedAt time.Time `json:"createdAt"`
}

// HomePageData is the aggregate payload for homepage render/API response.
type HomePageData struct {
	// Site contains singleton global homepage content.
	Site SiteContent `json:"site"`
	// Rooms contains room categories to render.
	Rooms []Room `json:"rooms"`
	// Gallery contains gallery cards to render.
	Gallery []GalleryItem `json:"gallery"`
	// Snapshot optionally includes last saved source snapshot.
	Snapshot *Snapshot `json:"snapshot,omitempty"`
}

// AdminPageData extends homepage payload with one-time UI flash messaging.
type AdminPageData struct {
	// HomePageData embeds base homepage content and collections.
	HomePageData
	// Flash is a short transient UI message for admin feedback.
	Flash string `json:"flash"`
}

// DefaultSiteContent returns starter content used during initial database seeding.
func DefaultSiteContent() SiteContent {
	return SiteContent{
		BrandName:    "Морской дворик",
		Tagline:      "Бутик-отдых у моря",
		HeroTitle:    "Уютный отель с двориком у моря в Феодосии",
		HeroSubtitle: "Комфортные номера, подогреваемый бассейн, просторная терраса и удобное прямое бронирование для летнего отдыха.",
		SeasonNote:   "Сезонная работа с 20 июня по 30 сентября: подогреваемый открытый бассейн и отдых в формате дворик-терраса.",
		BookingURL:   "#booking",
		Address:      "ул. Федько, 66А, Феодосия, Крым",
		Phone:        "+7 (978) 196-37-37",
		Email:        "hello@seaside-dvorik.local",
		CheckIn:      "14:00",
		CheckOut:     "12:00",
		AboutTitle:   "Гостеприимство с акцентом на комфорт",
		AboutBody:    "Этот шаблон повторяет структуру курортного сайта: выразительный первый экран, каталог номеров, блок удобств, галерея и контакты для бронирования. Перед публикацией замените тексты, фото и юридическую информацию на данные вашего объекта.",
	}
}

// DefaultRooms returns starter room categories used for first-run seeding.
func DefaultRooms() []Room {
	return []Room{
		{
			Name:        "Семейный номер на четверых",
			Summary:     "Светлый номер для семьи или небольшой компании с удобной планировкой и быстрым доступом к бассейну.",
			Capacity:    4,
			SizeSQM:     20,
			RoomCount:   1,
			FeatureNote: "Удобно для семейного отдыха",
			ImagePath:   "",
		},
		{
			Name:        "Улучшенный двухместный с балконом",
			Summary:     "Спокойный двухместный номер с зоной отдыха на балконе и комфортом в жаркий сезон.",
			Capacity:    2,
			SizeSQM:     17,
			RoomCount:   1,
			FeatureNote: "Собственный балкон",
			ImagePath:   "",
		},
		{
			Name:        "Мансардный люкс на четверых",
			Summary:     "Двухкомнатный вариант для длительного проживания: больше приватности и пространства.",
			Capacity:    4,
			SizeSQM:     38,
			RoomCount:   2,
			FeatureNote: "Двухкомнатный люкс",
			ImagePath:   "",
		},
	}
}

// DefaultGallery returns starter gallery entries used for first-run seeding.
func DefaultGallery() []GalleryItem {
	return []GalleryItem{
		{Title: "Главный фасад", Subtitle: "Первое впечатление со стороны улицы и уютный приватный двор.", Tone: "sunrise", ImagePath: ""},
		{Title: "Подогреваемый бассейн", Subtitle: "Курортная атмосфера для отдыха в жаркие летние дни.", Tone: "aqua", ImagePath: ""},
		{Title: "Терраса", Subtitle: "Открытая зона для утреннего кофе и вечернего отдыха.", Tone: "sand", ImagePath: ""},
		{Title: "Вид с балконов", Subtitle: "Небольшое личное пространство для гостей улучшенных номеров.", Tone: "twilight", ImagePath: ""},
	}
}
