import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label } from '@trussworks/react-uswds';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';

const AccountingCodes = ({ optional }) => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2>
          Accounting codes
          {optional && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>
        <Label>TAC</Label>
        <Label>SAC</Label>
      </Fieldset>
    </SectionWrapper>
  );
};

AccountingCodes.defaultProps = {
  optional: true,
  TACs: {},
  SACs: {},
};

AccountingCodes.propTypes = {
  optional: PropTypes.bool,
  TACs: PropTypes.shape({
    hhg: PropTypes.string,
    nts: PropTypes.string,
  }),
  SACs: PropTypes.shape({
    hhg: PropTypes.string,
    nts: PropTypes.string,
  }),
};

export default AccountingCodes;
