export function createItemRequest({ shipmentId, csrfToken, code, quantity1 }) {
  const tariff400ng = {
    '105B': { item: 'Pack Reg Crate', id: 'deb28967-d52c-4f04-8a0b-a264c9d80457', location: 'ORIGIN' },
    '105E': { item: 'UnPack Reg Crate', id: '6df4f1aa-a232-4eef-bbe8-f06bfb0b6d40', location: 'DESTINATION' },
    '35A': { item: 'Third Party Service', id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1', location: 'EITHER' },
    '226A': { item: 'Miscellaneous Charge', id: 'c5a6b126-de1a-4ab7-a5f4-c6e42cdf443b', location: 'EITHER' },
    '125A': {
      item: 'Shuttle Service 25 or less miles',
      id: 'e27968f9-44a7-4582-af85-f4a5891120fd',
      location: 'EITHER',
    },
    '125B': { item: 'Shuttle Service Over 25 Miles', id: 'd8044214-7fab-4588-8a6d-96d39d42e11d', location: 'EITHER' },
    '125C': {
      item: 'Shuttle Service 25 or less miles-OT',
      id: 'c85179e9-3986-463c-b970-d470133be993',
      location: 'EITHER',
    },
    '125D': {
      item: 'Shuttle Service Over 25 Miles-OT',
      id: '703c3e58-3835-439a-9ed1-6755eb91b62f',
      location: 'EITHER',
    },
  };
  let itemDetails;
  if (code in tariff400ng) {
    /* eslint-disable security/detect-object-injection */
    itemDetails = tariff400ng[code];
  } else {
    return;
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
