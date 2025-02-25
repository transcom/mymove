-- B-22653 Daniel Jordan update moving_expense_type to include SMALL_PACKAGE
ALTER TYPE moving_expense_type ADD VALUE IF NOT EXISTS 'SMALL_PACKAGE';