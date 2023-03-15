package router

type Router struct {
    routes []*Route
    middleware []Middleware
}
