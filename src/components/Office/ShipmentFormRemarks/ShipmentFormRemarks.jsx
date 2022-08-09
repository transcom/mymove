import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';
import { FormGroup } from '@material-ui/core';

import TextField from 'components/form/fields/TextField/TextField';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import formStyles from 'styles/form.module.scss';
import { officeRoles, roleTypes } from 'constants/userRoles';
import DataTable from 'components/DataTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint/index';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const ShipmentFormRemarks = ({ userRole, shipmentType, customerRemarks, counselorRemarks, error, showHint }) => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2 className={styles.SectionHeaderExtraSpacing}>
          Remarks{' '}
          {userRole === roleTypes.SERVICES_COUNSELOR && shipmentType !== SHIPMENT_OPTIONS.PPM && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>

        {shipmentType !== SHIPMENT_OPTIONS.PPM && (
          <DataTable custClass={styles.RemarksTable} columnHeaders={['Customer remarks']} dataRow={[customerRemarks]} />
        )}

        {userRole === roleTypes.TOO ? (
          <DataTable
            custClass={styles.RemarksTable}
            columnHeaders={['Counselor remarks']}
            dataRow={[counselorRemarks]}
          />
        ) : (
          <>
            {showHint && (
              <Hint>
                <p>500 characters</p>
              </Hint>
            )}
            <FormGroup className={styles.remarksField}>
              <TextField
                display="textarea"
                label="Counselor Remarks"
                data-testid="counselor-remarks"
                name="counselorRemarks"
                className={`${formStyles.remarks}`}
                placeholder=""
                id="counselorRemarks"
                maxLength={500}
                error={error}
              />
            </FormGroup>
          </>
        )}
      </Fieldset>
    </SectionWrapper>
  );
};

ShipmentFormRemarks.propTypes = {
  userRole: PropTypes.oneOf(officeRoles).isRequired,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
    SHIPMENT_OPTIONS.PPM,
  ]).isRequired,
  customerRemarks: PropTypes.string,
  counselorRemarks: PropTypes.string,
  showHint: PropTypes.bool,
  error: PropTypes.bool,
};

ShipmentFormRemarks.defaultProps = {
  customerRemarks: '—',
  counselorRemarks: '—',
  showHint: true,
  error: undefined,
};

export default ShipmentFormRemarks;
