package main

type Site struct {
	root string
}

func (s *Site) Enable() bool {
	return false
}

func (s *Site) Path() string {
	return "/my"
}

func (s *Site) Name() string {
	return "我的网站"
}

func (s *Site) Root() string {
	return s.root
}
