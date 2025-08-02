// Package service обеспечивает функциональность для управления маршрутами
// и предоставления к ним доступа.
package service

import (
	"github.com/ee-crocush/go-news/api-gateway/internal/infrastructure/config"
)

// RegistryService определяет интерфейс для доступа к зарегистрированным маршрутам.
type RegistryService interface {
	GetRouteByName(name string) (config.Route, bool)
	GetAllRoutes() []config.Route
}

// RouteRegistry хранит список маршрутов и предоставляет к ним доступ.
type RouteRegistry struct {
	routes []config.Route
}

// NewRouteRegistry создает новый экземпляр RouteRegistry.
// Принимает срез маршрутов и возвращает указатель на RouteRegistry.
func NewRouteRegistry(routes []config.Route) *RouteRegistry {
	return &RouteRegistry{routes: routes}
}

// GetRouteByName ищет маршрут по его имени.
// Принимает имя маршрута и возвращает найденный маршрут и булево значение, указывающее, был ли маршрут найден.
func (r *RouteRegistry) GetRouteByName(name string) (config.Route, bool) {
	for _, route := range r.routes {
		if route.Name == name {
			return route, true
		}
	}
	return config.Route{}, false
}

// GetAllRoutes возвращает все зарегистрированные маршруты.
func (r *RouteRegistry) GetAllRoutes() []config.Route {
	return r.routes
}
