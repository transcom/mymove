import { setFlagStyles, setDisplayFlags, getDisplayFlags, getMissingOrDash } from './displayFlags';

describe('setAndRetrieveFlags', () => {
  // example fields and data to use for testing, not reflective of reality
  const errorIfMissing = ['firstName'];
  const warnIfMissing = ['counselorRemarks'];
  const showWhenCollapsed = ['shipmentAddress'];
  const neverShow = ['postalCode'];
  const shipment = {
    lastName: 'LastName',
  };

  const styles = {
    row: 'row',
    warning: 'warning',
    missingInfoError: 'error',
  };

  it('can set and retrieve error flags', () => {
    setDisplayFlags(errorIfMissing, null, null, null, shipment);

    setFlagStyles(styles);

    const result = getDisplayFlags('firstName');

    expect(result.alwaysShow).toEqual(true);
    expect(result.classes).toEqual('row error');
  });

  it('can set and retrieve warning flags', () => {
    setDisplayFlags(null, warnIfMissing, null, null, shipment);

    setFlagStyles(styles);

    const result = getDisplayFlags('counselorRemarks');

    expect(result.alwaysShow).toEqual(true);
    expect(result.classes).toEqual('row warning');
  });

  it('can set and retrieve show when collapsed flags', () => {
    setDisplayFlags(null, null, showWhenCollapsed, null, shipment);

    setFlagStyles(styles);

    const result = getDisplayFlags('shipmentAddress');

    expect(result.alwaysShow).toEqual(true);
    expect(result.classes).toEqual('row');
    expect(result.hideRow).toEqual(false);
  });

  it('can set and retrieve never show flags', () => {
    setDisplayFlags(null, null, null, neverShow, shipment);

    setFlagStyles(styles);

    const result = getDisplayFlags('postalCode');

    expect(result.alwaysShow).toEqual(false);
    expect(result.hideRow).toEqual(true);
    expect(result.classes).toEqual('row');
  });

  it('will return missing or dash', () => {
    setDisplayFlags(errorIfMissing, null, null, null, shipment);

    const missing = getMissingOrDash('firstName');
    expect(missing).toEqual('Missing');

    const dash = getMissingOrDash('dashTest');
    expect(dash).toEqual('—');
  });
});
