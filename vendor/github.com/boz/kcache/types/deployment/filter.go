package deployment

import (
	"sort"

	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
)

func PodsFilter(sources ...*extv1beta1.Deployment) filter.ComparableFilter {

	// make a copy and sort
	srcs := make([]*extv1beta1.Deployment, len(sources))
	copy(srcs, sources)

	sort.Slice(srcs, func(i, j int) bool {
		if srcs[i].Namespace != srcs[j].Namespace {
			return srcs[i].Namespace < srcs[j].Namespace
		}
		return srcs[i].Name < srcs[j].Name
	})

	filters := make([]filter.Filter, 0, len(srcs))

	for _, svc := range srcs {

		var sfilter filter.Filter
		if sel := svc.Spec.Selector; sel != nil {
			sfilter = filter.LabelSelector(sel)
		} else {
			sfilter = filter.Labels(svc.Spec.Template.Labels)
		}

		nsfilter := filter.NSName(nsname.New(svc.GetNamespace(), ""))

		filters = append(filters, filter.And(nsfilter, sfilter))
	}

	return filter.Or(filters...)

}
