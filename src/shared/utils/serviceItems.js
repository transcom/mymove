import { dateSort } from 'shared/utils';

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
    if (shipmentServiceItems[`${shipmentId}`]) {
      return false;
    }

    shipmentServiceItems[`${shipmentId}`] = cards.filter((item) => item.shipmentId === shipmentId);
    shipmentOrder.push(serviceItem);
    return true;
  });

  shipmentOrder.sort(dateCreatedSort);

  const sortedCards = [];
  shipmentOrder.forEach((shipment) => {
    sortedCards.push(...shipmentServiceItems[`${shipment.shipmentId}`]);
    return true;
  });

  return sortedCards;
}
