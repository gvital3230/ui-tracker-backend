package ws

type VisitorId string
type ItemId string
type IsVisible bool
type DashBoardSessions map[VisitorId]map[ItemId]IsVisible
type Dashboard struct {
	ActiveSessions DashBoardSessions
}

func (d *Dashboard) Track(m TrackMessage) {
	//associate wsClient with browser visitor ID
	m.Client.VisitorId = m.Visitor

	if m.State {
		if _, ok := d.ActiveSessions[m.Visitor]; !ok {
			d.ActiveSessions[m.Visitor] = make(map[ItemId]IsVisible)
		}
		d.ActiveSessions[m.Visitor][m.ItemId] = true
	} else {
		if _, ok := d.ActiveSessions[m.Visitor][m.ItemId]; ok {
			delete(d.ActiveSessions[m.Visitor], m.ItemId)
		}
	}
}

func (d Dashboard) Unregister(visitor VisitorId) {
	if _, ok := d.ActiveSessions[visitor]; ok {
		delete(d.ActiveSessions, visitor)
	}
}
