class SurfBoard < ApplicationRecord
  belongs_to :user

  validates :name, :board_type, :length, :width, :thickness, :volume, :fin_setup, :brand, presence: true
  validates :length, :width, :thickness, :volume, numericality: { greater_than: 0 }
end
