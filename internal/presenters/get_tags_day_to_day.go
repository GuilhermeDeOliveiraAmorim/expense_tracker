package presenters

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/repositories"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/util"
)

type TagDayToDay struct {
	Amount float64 `json:"amount"`
	Name   string  `json:"name"`
	Color  string  `json:"color"`
}

type CategoryDayToDay struct {
	Amount float64       `json:"amount"`
	Name   string        `json:"name"`
	Color  string        `json:"color"`
	Tags   []TagDayToDay `json:"tags"`
}

type DayToDay struct {
	Date        time.Time          `json:"date"`
	NumberOfDay int                `json:"number_of_day"`
	DayOfWeek   string             `json:"day_of_week"`
	Month       string             `json:"month"`
	Year        string             `json:"year"`
	TotalAmount float64            `json:"total_amount"`
	Categories  []CategoryDayToDay `json:"categories"`
}

type GetTagsDayToDayInputDto struct {
	UserID string `json:"user_id"`
	Month  string `json:"month"`
	Year   string `json:"year"`
}

type GetTagsDayToDayOutputDto struct {
	Days []DayToDay `json:"days"`
}

type GetTagsDayToDayUseCase struct {
	PresentersRepository repositories.PresentersRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
}

func NewGetTagsDayToDayUseCase(
	PresentersRepository repositories.PresentersRepositoryInterface,
	UserRepository repositories.UserRepositoryInterface,
) *GetTagsDayToDayUseCase {
	return &GetTagsDayToDayUseCase{
		PresentersRepository: PresentersRepository,
		UserRepository:       UserRepository,
	}
}

func (c *GetTagsDayToDayUseCase) Execute(input GetTagsDayToDayInputDto) (GetTagsDayToDayOutputDto, []util.ProblemDetails) {
	user, err := c.UserRepository.GetUser(input.UserID)
	if err != nil {
		return GetTagsDayToDayOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Not Found",
				Title:    "User not found",
				Status:   404,
				Detail:   err.Error(),
				Instance: util.RFC404,
			},
		}
	} else if !user.Active {
		return GetTagsDayToDayOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Forbidden",
				Title:    "User is not active",
				Status:   403,
				Detail:   "User is not active",
				Instance: util.RFC403,
			},
		}
	}

	year, errYear := strconv.Atoi(input.Year)
	if errYear != nil {
		return GetTagsDayToDayOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid year",
				Status:   400,
				Detail:   errYear.Error(),
				Instance: util.RFC400,
			},
		}
	}

	month, errMonth := strconv.Atoi(input.Month)
	if errMonth != nil {
		return GetTagsDayToDayOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Bad Request",
				Title:    "Invalid month",
				Status:   400,
				Detail:   "Month is not in the correct format",
				Instance: util.RFC400,
			},
		}
	}

	expenses, err := c.PresentersRepository.GetTagsDayToDay(input.UserID, year, month)
	if err != nil {
		return GetTagsDayToDayOutputDto{}, []util.ProblemDetails{
			{
				Type:     "Internal Server Error",
				Title:    "Could not tags day by day",
				Status:   500,
				Detail:   err.Error(),
				Instance: util.RFC500,
			},
		}
	}

	daysMap := make(map[string]*DayToDay)

	for _, expense := range expenses {
		year := expense.ExpenseDate.Year()
		month := expense.ExpenseDate.Month().String()
		dayOfWeek := expense.ExpenseDate.Weekday().String()
		dateKey := expense.ExpenseDate.Format("2006-01-02")

		if _, exists := daysMap[dateKey]; !exists {
			daysMap[dateKey] = &DayToDay{
				Date:        expense.ExpenseDate,
				NumberOfDay: expense.ExpenseDate.Day(),
				DayOfWeek:   dayOfWeek,
				Month:       month,
				Year:        fmt.Sprint(year),
				TotalAmount: 0,
				Categories:  []CategoryDayToDay{},
			}
		}

		day := daysMap[dateKey]
		day.TotalAmount += expense.Amount

		categoryExists := false
		for i, category := range day.Categories {
			if category.Name == expense.Category.Name {
				categoryExists = true

				day.Categories[i].Amount += expense.Amount

				if len(expense.Tags) > 0 {
					day.Categories[i].Tags = append(day.Categories[i].Tags, TagDayToDay{
						Amount: expense.Amount,
						Name:   expense.Tags[0].Name,
						Color:  expense.Tags[0].Color,
					})
				}
				break
			}
		}

		if !categoryExists {
			category := CategoryDayToDay{
				Amount: expense.Amount,
				Name:   expense.Category.Name,
				Color:  expense.Category.Color,
			}

			if len(expense.Tags) > 0 {
				category.Tags = append(category.Tags, TagDayToDay{
					Amount: expense.Amount,
					Name:   expense.Tags[0].Name,
					Color:  expense.Tags[0].Color,
				})
			}

			day.Categories = append(day.Categories, category)
		}
	}

	var daysToDay []DayToDay
	for _, day := range daysMap {
		daysToDay = append(daysToDay, *day)
	}

	sort.Slice(daysToDay, func(i, j int) bool {
		return daysToDay[i].Date.Before(daysToDay[j].Date)
	})

	return GetTagsDayToDayOutputDto{
		Days: daysToDay,
	}, nil
}