package ws

type Dashboard struct {
	ActiveSessions map[string]map[string]bool
}

func (d *Dashboard) Track(m TrackMessage) {
	//associate wsClient with browser visitor ID
	m.Client.VisitorId = m.Visitor

	if m.State {
		if _, ok := d.ActiveSessions[m.Visitor]; !ok {
			d.ActiveSessions[m.Visitor] = make(map[string]bool)
		}
		d.ActiveSessions[m.Visitor][m.ItemId] = true
	} else {
		if _, ok := d.ActiveSessions[m.Visitor][m.ItemId]; ok {
			delete(d.ActiveSessions[m.Visitor], m.ItemId)
		}
	}
}

func (d Dashboard) Unregister(visitor string) {
	if _, ok := d.ActiveSessions[visitor]; ok {
		delete(d.ActiveSessions, visitor)
	}
}
