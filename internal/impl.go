package internal

import (
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/department"
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/department_role"
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/invitation"
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/member"
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/position"
	_ "github.com/LandcLi/LandcOrg/internal/dao_impl/position_role"

	_ "github.com/LandcLi/LandcOrg/internal/service_impl/department"
	_ "github.com/LandcLi/LandcOrg/internal/service_impl/invitation"
	_ "github.com/LandcLi/LandcOrg/internal/service_impl/member"
	_ "github.com/LandcLi/LandcOrg/internal/service_impl/org_role"
	_ "github.com/LandcLi/LandcOrg/internal/service_impl/position"

	_ "github.com/LandcLi/LandcOrg/internal/controller/department"
	_ "github.com/LandcLi/LandcOrg/internal/controller/invitation"
	_ "github.com/LandcLi/LandcOrg/internal/controller/member"
	_ "github.com/LandcLi/LandcOrg/internal/controller/org_role"
	_ "github.com/LandcLi/LandcOrg/internal/controller/position"
)
