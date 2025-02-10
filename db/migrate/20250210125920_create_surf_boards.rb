class CreateSurfBoards < ActiveRecord::Migration[8.0]
  def change
    create_table :surf_boards do |t|
      t.string  :name, null: false
      t.string  :board_type
      t.decimal :length, precision: 5, scale: 2
      t.decimal :width, precision: 5, scale: 2
      t.decimal :thickness, precision: 5, scale: 2
      t.decimal :volume, precision: 5, scale: 2
      t.string  :fin_setup
      t.string  :brand
      t.text    :notes
      t.references :user, null: false, foreign_key: true

      t.timestamps
    end
  end
end
