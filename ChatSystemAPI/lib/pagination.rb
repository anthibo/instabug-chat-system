module Pagination
  def paginate(page: 1, per_page: 10)
    offset((page - 1) * per_page).limit(per_page)
  end
end
