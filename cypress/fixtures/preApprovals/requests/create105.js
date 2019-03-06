export function createItemRequest({ shipmentId, csrfToken, code, quantity1 }) {
  const tariff400ng = {
    '105B': { item: 'Pack Reg Crate', id: 'deb28967-d52c-4f04-8a0b-a264c9d80457', location: 'ORIGIN' },
    '105E': { item: 'UnPack Reg Crate', id: '6df4f1aa-a232-4eef-bbe8-f06bfb0b6d40', location: 'DESTINATION' },
  };
  let itemDetails;
  if (code === '105B') {
    itemDetails = tariff400ng['105B'];
  } else if (code === '105E') {
    itemDetails = tariff400ng['105E'];
  } else {
    itemDetails = { item: 'unknown', id: 'unknown', location: 'BOTH' };
  }

  return {
    method: 'POST',
    url: `/api/v1/shipments/${shipmentId}/accessorials`,
    headers: {
      'X-CSRF-TOKEN': csrfToken,
    },
    body: {
      tariff400ng_item: {
        code: code,
        created_at: '2019-03-05T15:34:29.785Z',
        discount_type: 'HHG',
        id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
        item: itemDetails.item,
        location: itemDetails.location,
        ref_code: 'NONE',
        requires_pre_approval: true,
        uom_1: 'CF',
        uom_2: 'NONE',
        updated_at: '2019-03-05T15:34:29.785Z',
      },
      location: 'ORIGIN',
      quantity_1: quantity1 * 10000,
      notes: `notes notes ${code}`,
      tariff400ng_item_id: itemDetails.id,
    },
  };
}
