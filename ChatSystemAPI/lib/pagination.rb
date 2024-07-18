module Pagination
  def paginate(relation, page: 1, per_page: 10)
    relation.offset((page - 1) * per_page).limit(per_page)
  end
end
