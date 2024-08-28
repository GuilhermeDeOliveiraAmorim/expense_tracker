package handlers

import (
	"net/http"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/infra/factory"
	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/presenters"
	"github.com/gin-gonic/gin"
)

type PresentersHandler struct {
	presenterFactory *factory.PresentersFactory
}

func NewPresentersHandler(factory *factory.PresentersFactory) *PresentersHandler {
	return &PresentersHandler{
		presenterFactory: factory,
	}
}

func (p *PresentersHandler) ShowTotalExpensesCategoryPeriod(c *gin.Context) {
	var input presenters.ShowTotalExpensesCategoryPeriodInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := p.presenterFactory.ShowTotalExpensesCategoryPeriod.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.JSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.JSON(err.Status, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(http.StatusOK, output)
}

func (p *PresentersHandler) ShowCategoryTreemapAmountPeriod(c *gin.Context) {
	var input presenters.ShowCategoryTreemapAmountPeriodInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := p.presenterFactory.ShowCategoryTreemapAmountPeriod.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.JSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.JSON(err.Status, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(http.StatusOK, output)
}

func (p *PresentersHandler) ShowExpenseSimpleTablePeriod(c *gin.Context) {
	var input presenters.ShowExpenseSimpleTablePeriodInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, erros := p.presenterFactory.ShowExpenseSimpleTablePeriod.Execute(input)
	if len(erros) > 0 {
		for _, err := range erros {
			if err.Status == 500 {
				c.JSON(err.Status, gin.H{"error": err})
				return
			} else {
				c.JSON(err.Status, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(http.StatusOK, output)
}
