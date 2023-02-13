package query

// Return instance of Query builder with default query
type QueryBuilderFactory interface {
	Create(query string) QueryBuilder
}

type queryBuilderFactory struct {
}

func NewQueryBuilderFactory() QueryBuilderFactory {
	return queryBuilderFactory{}
}

func (qbf queryBuilderFactory) Create(query string) QueryBuilder {
	return &psqlQueryBuilder{baseQuery: query, expressions: make([]Expression, 0)}
}

// Resolve query
type QueryBuilder interface {
	AddPaging(limit, offset int)
	AddSorting(sorts []string)
	GetArguments() []any
	GetQuery() string
	ILike(propertyName string, value any, matchMode int)
	Eq(propertyName string, value any)
	In(propertyName string, value []any)
	Exists(query string)
	NotExists(query string)
	Condition() Expression
	// Condition(expr Expression) Expression// create condition
	/*
								expression = builder.condition().eq(
							                            fieldName,
							                            filterItem.getValue().stream().map(BooleanUtils::toBoolean).findFirst().orElse(null)
							                    );
								expression = builder.condition().or(
						                            filterItem.getValue().stream().map(
						                                    itemValue -> filterItem.isLike() ?
						                                            builder.condition().ilike(fieldName, itemValue, MatchMode.ANYWHERE) :
						                                            builder.condition().eq(fieldName, itemValue)
						                            ).toArray(Expression[]::new)
								builder.condition().parentheses(expression)



								private Expression createHasNotActualDeedExpression(ExtQueryBuilder builder){
				        Expression deedStatusIdIsNull = builder.condition().isNull("deed.status_id");
				        Expression deedStatusClosed = builder.condition().eq("deed.status_id", DeedStatusEnum.CLOSED.getId());
				        Expression deedNotExistsOrClosed = builder.condition().parentheses(builder.condition().or(deedStatusIdIsNull, deedStatusClosed));
				        Expression applicantNotExistsInOtherNotClosedDeed = builder.condition().not(new ExistsExpression(
				                queryBuilderFactory.create("SELECT 1 FROM housing_subsidies.deed d1 WHERE d1.applicant_id = individual.id ").
				                        not(builder.condition().eq("d1.status_id", DeedStatusEnum.CLOSED.getId()))));
				        return builder.condition().and(deedNotExistsOrClosed, applicantNotExistsInOtherNotClosedDeed);
				    }




					QueryBuilder<?> builder = queryBuilderFactory.create("SELECT " +
			                "work.id, work.individual_person_id, work.employer_id, work.work_func, work.employment_date, " +
			                "work.dismissal_date, work.dismissal_reason, work.base_doc, work.is_pfr, work.create_user, " +
			                "work.create_date, work.change_user, work.change_date, work.check_date, work.check_user, " +
			                "work.employer_id, " +
			                "e.short_name AS employer_short_name " +
			                "FROM housing_subsidies.work as work " +
			                "LEFT JOIN housing_subsidies.employer e ON work.employer_id = e.id ");
			        builder.in("work.individual_person_id", individualIds.toArray(Number[]::new));
			        builder.in("work.employer_id", employerIds.toArray(Number[]::new));
			        builder.isNull("work.dismissal_date");
			        List<Long> commonList = new ArrayList<>(individualIds);
			        commonList.addAll(employerIds);
			        return jdbcTemplate.query(builder.getQuery(), rowMapper, commonList.toArray());



					if (BooleanUtils.isTrue(query.isCertificated())) {
		            builder.exists("SELECT 1 FROM housing_subsidies.certificate cert WHERE cert.deed_id = plm.deed_id");
		        } else {
		            builder.notExists("SELECT 1 FROM housing_subsidies.certificate cert WHERE cert.deed_id = plm.deed_id");
		        }
	*/
}

type psqlQueryBuilder struct {
	baseQuery   string
	expressions []Expression
	limit       int
	offset      int
	values      map[int]any
}

func (qb *psqlQueryBuilder) AddPaging(limit, offset int) {
	qb.limit = limit
	qb.offset = offset
}

func (qb *psqlQueryBuilder) AddSorting(sorts []string) {

}

func (qb *psqlQueryBuilder) GetArguments() []any {

	return make([]any, 0)
}

func (qb *psqlQueryBuilder) GetQuery() string {

	return ""
}

func (qb *psqlQueryBuilder) ILike(propertyName string, value any, matchMode int) {
	qb.expressions = append(qb.expressions)
}

func (qb *psqlQueryBuilder) Eq(propertyName string, value any) {

}

func (qb *psqlQueryBuilder) In(propertyName string, value []any) {

}

func (qb *psqlQueryBuilder) Exists(query string) {

}

func (qb *psqlQueryBuilder) NotExists(query string) {

}

func (qb *psqlQueryBuilder) Condition() Expression {
	return nil
}
