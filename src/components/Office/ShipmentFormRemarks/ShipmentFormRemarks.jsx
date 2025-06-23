import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, FormGroup } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import formStyles from 'styles/form.module.scss';
import { officeRoles, roleTypes } from 'constants/userRoles';
import DataTable from 'components/DataTable';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import Hint from 'components/Hint/index';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const ShipmentFormRemarks = ({
  userRole,
  shipmentType,
  customerRemarks,
  counselorRemarks,
  error,
  showHint,
  advanceStatus,
}) => {
  const advanceRejected = advanceStatus === 'REJECTED';
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2 className={styles.SectionHeaderExtraSpacing}>
          Remarks{' '}
          {userRole === roleTypes.SERVICES_COUNSELOR && shipmentType !== SHIPMENT_OPTIONS.PPM && (
            <span className="float-right">
              <span className={formStyles.optional} />
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
              {advanceRejected && requiredAsteriskMessage}
              <TextField
                display="textarea"
                label="Counselor remarks"
                data-testid="counselor-remarks"
                name="counselorRemarks"
                className={`${formStyles.remarks}`}
                placeholder=""
                id="counselorRemarks"
                maxLength={500}
                error={error}
                showRequiredAsterisk={advanceRejected}
                required={advanceRejected}
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
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_OPTIONS.MOBILE_HOME,
    SHIPMENT_OPTIONS.BOAT,
    SHIPMENT_TYPES.BOAT_HAUL_AWAY,
    SHIPMENT_TYPES.BOAT_TOW_AWAY,
    SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
  ]).isRequired,
  customerRemarks: PropTypes.string,
  counselorRemarks: PropTypes.string,
  showHint: PropTypes.bool,
  error: PropTypes.bool,
  advanceStatus: PropTypes.string,
};

ShipmentFormRemarks.defaultProps = {
  customerRemarks: '—',
  counselorRemarks: '—',
  showHint: true,
  error: undefined,
  advanceStatus: undefined,
};

export default ShipmentFormRemarks;
