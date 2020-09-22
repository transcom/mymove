import { dateSort } from 'shared/utils';
import { convertFromThousandthInchToInch } from 'shared/formatters';

export function sortServiceItemsByGroup(serviceItemCards) {
  // Make a copy so we're not mutating the props
  const cards = [...serviceItemCards];
  // Will populate with earliest service item of each shipment id
  const shipmentOrder = [];
  // Contains sorted service items keyed by shipment id or undefined for basic items
  const shipmentServiceItems = {};

  const dateCreatedSort = dateSort('createdAt', 'asc');

  cards.sort(dateCreatedSort);

  cards.forEach((serviceItem) => {
    const { shipmentId } = serviceItem;
    // We've already added the earliest service item for this shipment, continue until we get to the next
    if (!shipmentServiceItems[`${shipmentId}`]) {
      shipmentServiceItems[`${shipmentId}`] = cards.filter((item) => item.shipmentId === shipmentId);
      shipmentOrder.push(serviceItem);
    }
  });

  shipmentOrder.sort(dateCreatedSort);

  const sortedCards = [];
  shipmentOrder.forEach((shipment) => {
    sortedCards.push(...shipmentServiceItems[`${shipment.shipmentId}`]);
  });

  return sortedCards;
}

export function formatDimensions(dimensions, conversion = convertFromThousandthInchToInch, symbol = '"') {
  if (!dimensions) {
    return '';
  }

  return `${conversion(dimensions.length)}${symbol}x${conversion(dimensions.width)}${symbol}x${conversion(
    dimensions.height,
  )}${symbol}`;
}
