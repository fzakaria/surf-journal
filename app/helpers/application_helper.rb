module ApplicationHelper
  def active_nav_class(path)
    if current_page?(path)
      "text-primary"
    else
      ""
    end
  end
end
