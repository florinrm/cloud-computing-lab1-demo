package gateway

import (
	"context"
	"encoding/json"
	"exercise2/domain"
	"exercise2/repository"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	bookPath = "books"
)

type API struct {
	repo *repository.BookRepository
}

func NewAPI(repo *repository.BookRepository) *API {
	return &API{
		repo: repo,
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/app")
	ws.Route(ws.POST(bookPath).
		To(api.addBookHandler).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Writes back a json with what you gave it"))
	ws.Route(ws.GET(bookPath).
		To(api.getBooksHandler).
		Produces(restful.MIME_JSON).
		Doc("Writes back a json with what you gave it"))
}

func (api *API) addBookHandler(req *restful.Request, resp *restful.Response) {
	ctx := context.Background()
	body := req.Request.Body
	if body == nil {
		log.WithContext(ctx).Errorf("Couldn't read request body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body"))
		return

	}
	defer func() {
		_ = body.Close()
	}()

	var err error
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't read request body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	data, err := io.ReadAll(body)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't read request body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	book := &domain.Book{}
	err = json.Unmarshal(data, book)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't unmarshal request body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	book, err = api.repo.AddBook(ctx, book)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't insert book")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	err = resp.WriteAsJson(map[string]*domain.Book{"res": book})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't write response body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
}

func (api *API) getBooksHandler(req *restful.Request, resp *restful.Response) {
	ctx := context.Background()
	books, err := api.repo.GetBooks(ctx)

	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to retrieve books")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	if books == nil {
		books = make([]domain.Book, 0)
	}

	err = resp.WriteAsJson(map[string][]domain.Book{
		"result": books,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Couldn't write response body")
		_ = resp.WriteServiceError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
}
