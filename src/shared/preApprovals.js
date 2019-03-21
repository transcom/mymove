export const isNewAccessorial = item => {
  if (!item) return false;

  const code = item.tariff400ng_item.code;
  if ((code === '105B' || code === '105E') && !item.crate_dimensions) {
    return false;
  }
  if (code === '35A' && !item.estimate_amount_cents) {
    return false;
  }

  return true;
};
