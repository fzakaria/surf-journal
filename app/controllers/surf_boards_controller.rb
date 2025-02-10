class SurfBoardsController < ApplicationController
  before_action :set_surf_board, only: [:show, :edit, :update, :destroy]

  def index
    @surf_boards = SurfBoard.all
  end

  def show
  end

  def new
    @surf_board = SurfBoard.new
  end

  def create
    @surf_board = SurfBoard.new(surf_board_params)
    @surf_board.user = Current.user

    if @surf_board.save
      redirect_to @surf_board, notice: 'Surfboard was successfully created.'
    else
      render :new, status: :unprocessable_entity
    end
  end

  def edit
  end

  def update
    if @surf_board.update(surf_board_params)
      redirect_to @surf_board, notice: 'Surfboard was successfully updated.'
    else
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @surf_board.destroy
    redirect_to surf_boards_path, notice: 'Surfboard was successfully deleted.'
  end

  private

  def set_surf_board
    @surf_board = SurfBoard.find(params[:id])
  end

  def surf_board_params
    params.require(:surf_board).permit(:name, :board_type, :length, :width, :thickness, :volume, :fin_setup, :brand, :notes)
  end
end
