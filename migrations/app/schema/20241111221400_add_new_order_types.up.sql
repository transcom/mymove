-- adding order types to orders_type enum used in the orders table.
ALTER TYPE public.orders_type ADD VALUE 'EARLY_RETURN_OF_DEPENDENTS';
ALTER TYPE public.orders_type ADD VALUE 'STUDENT_TRAVEL';
COMMENT ON COLUMN orders.orders_type IS 'MilMove supports 10 orders types: Permanent change of station (PCS), local move, retirement, separation, wounded warrior, bluebark, safety, temporary duty (TDY), early return of dependents, and student travel.
In general, the moving process starts with the job/travel orders a customer receives from their service. In the orders, information describing rank, the duration of job/training, and their assigned location will determine if their entire dependent family can come, what the customer is allowed to bring, and how those items will arrive to their new location.'
