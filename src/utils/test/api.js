import { INACCESSIBLE_API_RESPONSE } from 'shared/Inaccessible';

// Expected return from a query when an office user tries to access a safety move without the adequate permissions
export const INACCESSIBLE_RETURN_VALUE = {
  isLoading: false,
  isError: true,
  isSuccess: false,
  errors: [{ response: { body: { message: INACCESSIBLE_API_RESPONSE } } }],
};

// Expected generic response for server side errors
export const ERROR_RETURN_VALUE = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

// Expected response when a query is still running on the backend
export const LOADING_RETURN_VALUE = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};
