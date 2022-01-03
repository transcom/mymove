import React, { useState } from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';

const ShipmentVendor = () => {
  const [inputProps, , helperProps] = useField('usesExternalVendor');
  const [hasChangedValue, setHasChangedValue] = useState(false);

  const usingExternalVendor = inputProps.value === true;

  const handleChangeToPrime = () => {
    helperProps.setValue(false);
    setHasChangedValue(true);
  };
  const handleChangeToExternal = () => {
    helperProps.setValue(true);
    setHasChangedValue(true);
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Vendor</h2>

        <Grid row gap>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label}>Who will handle this shipment?</Label>
              <Radio
                id="vendorPrime"
                label="GHC prime contractor"
                name="usesExternalVendor"
                value="GHC"
                title="GHC prime contractor"
                checked={!usingExternalVendor}
                onChange={handleChangeToPrime}
              />
              <Radio
                id="vendorExternal"
                label="External vendor"
                name="usesExternalVendor"
                value="External"
                title="External vendor"
                checked={usingExternalVendor}
                onChange={handleChangeToExternal}
              />
            </FormGroup>

            {(usingExternalVendor || hasChangedValue) && (
              <div className={styles.VendorExplainerText}>
                {usingExternalVendor ? (
                  <ul>
                    <li>This shipment won&apos;t be sent to the GHC prime contractor.</li>
                    <li>Shipment details will not automatically be shared with the movers handling it.</li>
                  </ul>
                ) : (
                  <>This shipment will be sent to the GHC prime contractor.</>
                )}
              </div>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentVendor;
