export function selectCustomer(state, customerId) {
  if (!state.entities.customer) {
    return {};
  }

  return Object.values(state.entities.customer).filter(customer => customerId === customer.id)[0] || {};
}
