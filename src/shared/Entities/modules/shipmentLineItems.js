import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { shipmentLineItems as ShipmentLineItemsModel } from '../schema';
import { denormalize } from 'normalizr';
import { get, orderBy, filter, map, keys } from 'lodash';
import { createSelector } from 'reselect';

export function createShipmentLineItem(label, shipmentId, payload) {
  return swaggerRequest(getPublicClient, 'accessorials.createShipmentLineItem', { shipmentId, payload }, { label });
}

export function updateShipmentLineItem(label, shipmentLineItemId, payload) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.updateShipmentLineItem',
    { shipmentLineItemId, payload },
    { label },
  );
}

export function deleteShipmentLineItem(label, shipmentLineItemId) {
  return swaggerRequest(getPublicClient, 'accessorials.deleteShipmentLineItem', { shipmentLineItemId }, { label });
}

export function approveShipmentLineItem(label, shipmentLineItemId) {
  return swaggerRequest(getPublicClient, 'accessorials.approveShipmentLineItem', { shipmentLineItemId }, { label });
}

export function getAllShipmentLineItems(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'accessorials.getShipmentLineItems', { shipmentId }, { label });
}

const selectShipmentLineItems = state => {
  return denormalize(keys(get(state, 'entities.shipmentLineItems', {})), ShipmentLineItemsModel, state.entities);
};

export const selectSortedShipmentLineItems = createSelector([selectShipmentLineItems], shipmentLineItems =>
  orderBy(shipmentLineItems, ['status', 'approved_date', 'submitted_date'], ['asc', 'desc', 'desc']),
);

export const selectSortedPreApprovalShipmentLineItems = createSelector(
  [selectSortedShipmentLineItems],
  shipmentLineItems => filter(shipmentLineItems, lineItem => lineItem.tariff400ng_item.requires_pre_approval),
);

export const getShipmentLineItemsLabel = 'ShipmentLineItems.getAllShipmentLineItems';
export const createShipmentLineItemLabel = 'ShipmentLineItems.createShipmentLineItem';
export const deleteShipmentLineItemLabel = 'ShipmentLineItems.deleteShipmentLineItem';
export const approveShipmentLineItemLabel = 'ShipmentLineItems.approveShipmentLineItem';
export const updateShipmentLineItemLabel = 'ShipmentLineItems.updateShipmentLineItem';

export const selectShipmentLineItem = (state, id) => denormalize([id], ShipmentLineItemsModel, state.entities)[0];

const getUnbilledShipmentLineItems = (state, shipmentId) => {
  const items = filter(state.entities.shipmentLineItems, item => {
    return item.shipment_id === shipmentId && !item.invoice_id;
  });

  //this denormalize step can be skipped because tariff400ng_item data is already available under items
  //but this is the right way to hydrate the data structure so leaving it in
  let denormItems = denormalize(map(items, 'id'), ShipmentLineItemsModel, state.entities);
  return denormItems.filter(item => {
    return !item.tariff400ng_item.requires_pre_approval || item.status === 'APPROVED';
  });
};

export const makeGetUnbilledShipmentLineItems = () => createSelector([getUnbilledShipmentLineItems], items => items);

export const makeTotalFromUnbilledLineItems = () =>
  createSelector([getUnbilledShipmentLineItems], items => {
    return items.reduce((acm, item) => {
      return acm + item.amount_cents;
    }, 0);
  });
