// Code generated by entc, DO NOT EDIT.

package ent

import "context"

func (c *Competition) Recruitments(ctx context.Context) ([]*Recruitment, error) {
	result, err := c.Edges.RecruitmentsOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryRecruitments().All(ctx)
	}
	return result, err
}

func (pr *Prefecture) Recruitments(ctx context.Context) ([]*Recruitment, error) {
	result, err := pr.Edges.RecruitmentsOrErr()
	if IsNotLoaded(err) {
		result, err = pr.QueryRecruitments().All(ctx)
	}
	return result, err
}

func (r *Recruitment) User(ctx context.Context) (*User, error) {
	result, err := r.Edges.UserOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryUser().Only(ctx)
	}
	return result, err
}

func (r *Recruitment) Prefecture(ctx context.Context) (*Prefecture, error) {
	result, err := r.Edges.PrefectureOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryPrefecture().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (r *Recruitment) Competition(ctx context.Context) (*Competition, error) {
	result, err := r.Edges.CompetitionOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryCompetition().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (u *User) Recruitments(ctx context.Context) ([]*Recruitment, error) {
	result, err := u.Edges.RecruitmentsOrErr()
	if IsNotLoaded(err) {
		result, err = u.QueryRecruitments().All(ctx)
	}
	return result, err
}
