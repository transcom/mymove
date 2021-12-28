import React from 'react';
import PropTypes from 'prop-types';
import { Label, Fieldset, Textarea } from '@trussworks/react-uswds';
import { Field } from 'formik';

import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import formStyles from 'styles/form.module.scss';
import { officeRoles, roleTypes } from 'constants/userRoles';
import DataTable from 'components/DataTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint/index';

const ShipmentFormRemarks = ({ userRole, customerRemarks, counselorRemarks }) => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2 className={styles.SectionHeaderExtraSpacing}>
          Remarks{' '}
          {userRole === roleTypes.SERVICES_COUNSELOR && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>

        <DataTable custClass={styles.RemarksTable} columnHeaders={['Customer remarks']} dataRow={[customerRemarks]} />

        {userRole === roleTypes.TOO ? (
          <DataTable
            custClass={styles.RemarksTable}
            columnHeaders={['Counselor remarks']}
            dataRow={[counselorRemarks]}
          />
        ) : (
          <>
            <Label htmlFor="counselorRemarks">Counselor remarks</Label>
            <Hint>
              <p>500 characters</p>
            </Hint>
            <Field
              as={Textarea}
              data-testid="counselor-remarks"
              name="counselorRemarks"
              className={`${formStyles.remarks}`}
              placeholder=""
              id="counselorRemarks"
              maxLength={500}
            />
          </>
        )}
      </Fieldset>
    </SectionWrapper>
  );
};

ShipmentFormRemarks.propTypes = {
  userRole: PropTypes.oneOf(officeRoles).isRequired,
  customerRemarks: PropTypes.string,
  counselorRemarks: PropTypes.string,
};

ShipmentFormRemarks.defaultProps = {
  customerRemarks: '—',
  counselorRemarks: '—',
};

export default ShipmentFormRemarks;
