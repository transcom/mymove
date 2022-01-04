import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Button } from '@trussworks/react-uswds';

import styles from './ShipmentAccountingCodes.module.scss';

import formStyles from 'styles/form.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import AccountingCodeSection from 'components/Office/AccountingCodeSection/AccountingCodeSection';
import { AccountingCodesShape } from 'types/accountingCodes';

const AccountingCodes = ({ optional, TACs, SACs, onEditCodesClick }) => {
  const hasCodes = Object.keys(TACs).length + Object.keys(SACs).length > 0;

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2 className={shipmentFormStyles.SectionHeaderExtraSpacing}>
          Accounting codes
          {optional && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>

        <AccountingCodeSection
          label="TAC"
          emptyMessage="No TAC code entered."
          fieldName="tacType"
          shipmentTypes={TACs}
        />
        <AccountingCodeSection
          label="SAC (optional)"
          emptyMessage="No SAC code entered."
          fieldName="sacType"
          shipmentTypes={SACs}
        />

        <Button type="button" onClick={onEditCodesClick} secondary={hasCodes} className={styles.addCodeBtn}>
          {hasCodes ? 'Add or edit codes' : 'Add code'}
        </Button>
      </Fieldset>
    </SectionWrapper>
  );
};

AccountingCodes.defaultProps = {
  optional: true,
  TACs: {},
  SACs: {},
  onEditCodesClick: () => {},
};

AccountingCodes.propTypes = {
  optional: PropTypes.bool,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  onEditCodesClick: PropTypes.func,
};

export default AccountingCodes;
