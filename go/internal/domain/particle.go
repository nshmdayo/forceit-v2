package domain

type Particle struct {
	Position Vec3
	Velocity Vec3
	Radius   float64
}

func (p *Particle) Integrate(dt float64) {
	p.Position = p.Position.Add(p.Velocity.Scale(dt))
}
