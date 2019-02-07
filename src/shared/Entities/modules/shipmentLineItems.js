import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { shipmentLineItems as ShipmentLineItemsModel } from '../schema';
import { denormalize } from 'normalizr';
import { get, orderBy, filter, map, keys, flow } from 'lodash';
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

// Show linehaul (and related) items before any accessorial items by adding isLinehaul property.
function listLinehaulItemsBeforeAccessorials(items) {
  const linehaulRelatedItems = ['LHS', '135A', '135B', '105A', '105C', '16A'];
  return items.map(item => {
    return {
      ...item,
      isLinehaul: linehaulRelatedItems.includes(item.tariff400ng_item.code) ? 1 : 10,
    };
  });
}

function orderItemsBy(items) {
  const sortOrder = {
    fields: ['isLinehaul', 'status', 'approved_date', 'submitted_date', 'tariff400ng_item.code'],
    order: ['asc', 'asc', 'desc', 'desc', 'desc'],
  };
  return orderBy(items, sortOrder.fields, sortOrder.order);
}

const selectShipmentLineItems = (state, shipmentId) => {
  let filteredItems = denormalize(
    keys(get(state, 'entities.shipmentLineItems', {})),
    ShipmentLineItemsModel,
    state.entities,
  );
  //only filter by shipmentId if it is explicitly passed
  if (!shipmentId) {
    return filteredItems;
  }
  return filter(filteredItems, item => {
    return item.shipment_id === shipmentId;
  });
};

export const selectSortedShipmentLineItems = createSelector([selectShipmentLineItems], items =>
  flow([listLinehaulItemsBeforeAccessorials, orderItemsBy])(items),
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

const selectInvoicesShipmentLineItemsByInvoiceId = (state, invoiceId) => {
  const items = filter(state.entities.shipmentLineItems, item => {
    return item.invoice_id === invoiceId;
  });

  return denormalize(map(items, 'id'), ShipmentLineItemsModel, state.entities);
};

const selectUnbilledShipmentLineItemsByShipmentId = (state, shipmentId) => {
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

export const selectUnbilledShipmentLineItems = createSelector([selectUnbilledShipmentLineItemsByShipmentId], items =>
  flow([listLinehaulItemsBeforeAccessorials, orderItemsBy])(items),
);

export const selectInvoiceShipmentLineItems = createSelector([selectInvoicesShipmentLineItemsByInvoiceId], items =>
  flow([listLinehaulItemsBeforeAccessorials, orderItemsBy])(items),
);

export const selectTotalFromUnbilledLineItems = createSelector([selectUnbilledShipmentLineItemsByShipmentId], items => {
  return items.reduce((acm, item) => {
    return acm + item.amount_cents;
  }, 0);
});

export const selectTotalFromInvoicedLineItems = createSelector([selectInvoicesShipmentLineItemsByInvoiceId], items => {
  return items.reduce((acm, item) => {
    return acm + item.amount_cents;
  }, 0);
});

export const selectLocationFromTariff400ngItem = (state, selectedTariff400ngItem) => {
  if (!selectedTariff400ngItem) return [];
  const lineItemLocations = get(state, 'swaggerPublic.spec.definitions.ShipmentLineItem', {}).properties.location;
  if (!lineItemLocations.enum) return [];
  const tariff400ngItemLocation = selectedTariff400ngItem.location;
  // Choose location options based on tariff400ng choice.
  return lineItemLocations.enum.filter(lineItemLocation => {
    return tariff400ngItemLocation === 'EITHER'
      ? lineItemLocation === 'ORIGIN' || lineItemLocation === 'DESTINATION'
      : lineItemLocation === tariff400ngItemLocation;
  });
};
