import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getPublicClient } from 'shared/Swagger/api';
import { invoice as InvoiceModel, invoices as InvoicesModel } from '../schema';
import { denormalize } from 'normalizr';
import { get, orderBy, filter, keys } from 'lodash';
import { createSelector } from 'reselect';

export const getShipmentInvoicesLabel = 'Shipments.getShipmentInvoices';
export const createInvoiceLabel = 'Shipments.createAndSendHHGInvoice';

export function createInvoice(label, shipmentId) {
  return swaggerRequest(getClient, 'shipments.createAndSendHHGInvoice', { shipmentId }, { label });
}

export function getAllInvoices(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'shipments.getShipmentInvoices', { shipmentId }, { label });
}

const selectInvoices = (state, shipmentId) => {
  let filteredItems = denormalize(keys(get(state, 'entities.invoices', {})), InvoicesModel, state.entities);
  //only filter by shipmentId if it is explicitly passed
  if (!shipmentId) {
    return filteredItems;
  }
  return filter(filteredItems, item => {
    return item.shipment_id === shipmentId;
  });
};

export const selectSortedInvoices = createSelector([selectInvoices], items =>
  orderBy(items, ['invoiced_date'], ['desc']),
);

export const selectInvoice = (state, id) => denormalize([id], InvoiceModel, state.entities)[0];
