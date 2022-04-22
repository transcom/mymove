import { dateSort } from 'shared/utils';
import { convertFromThousandthInchToInch } from 'utils/formatters';

function sortServiceItemsIntoShipments(serviceItems) {
  // Make a copy so we're not mutating the props
  const cards = [...serviceItems];
  // Will populate with earliest service item of each shipment id
  const shipmentOrder = [];
  // Contains sorted service items keyed by shipment id or undefined for basic items
  const shipmentServiceItems = {};

  const dateCreatedSort = dateSort('createdAt', 'asc');

  cards.sort(dateCreatedSort);

  cards.forEach((serviceItem) => {
    const { mtoShipmentID } = serviceItem;
    // We've already added the earliest service item for this shipment, continue until we get to the next
    if (!shipmentServiceItems[`${mtoShipmentID}`]) {
      shipmentServiceItems[`${mtoShipmentID}`] = cards.filter((item) => item.mtoShipmentID === mtoShipmentID);
      shipmentOrder.push(serviceItem);
    }
  });

  shipmentOrder.sort(dateCreatedSort);

  const sortedShipments = [];
  shipmentOrder.forEach((shipment) => {
    sortedShipments.push(shipmentServiceItems[`${shipment.mtoShipmentID}`]);
  });

  return sortedShipments;
}

// Returns an array of service item objects grouped by shipment, sorted by creation date
export function sortServiceItemsByGroup(serviceItemCards) {
  const sortedShipments = sortServiceItemsIntoShipments(serviceItemCards);

  const sortedCards = [];
  sortedShipments.forEach((shipment) => {
    sortedCards.push(...shipment);
  });

  return sortedCards;
}

// Sorts an array of service items into grouped shipments (basic, hhg-shipment1, hhg-shimpent2, ntsr, ...)
export function groupByShipment(serviceItems) {
  return sortServiceItemsIntoShipments(serviceItems);
}

export function formatDimensions(dimensions, conversion = convertFromThousandthInchToInch, symbol = '"') {
  if (!dimensions) {
    return '';
  }

  return `${conversion(dimensions.length)}${symbol}x${conversion(dimensions.width)}${symbol}x${conversion(
    dimensions.height,
  )}${symbol}`;
}
