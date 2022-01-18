import React from 'react';
import { string } from 'prop-types';
import { FormGroup, Label, TextInput } from '@trussworks/react-uswds';

import styles from './PrintableLegalese.module.scss';

import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import CertificationText from 'scenes/Legalese/CertificationText';
import { formatSignatureDate } from 'utils/formatters';

const PrintableLegalese = ({ signature, signatureDate }) => (
  <div className={styles.printableCertification}>
    <CertificationText certificationText={completeCertificationText} />
    <div className={styles.signatureFields}>
      <FormGroup className={styles.signature}>
        <Label htmlFor="signature">Signature</Label>
        <TextInput name="signature" readOnly value={signature} />
      </FormGroup>
      <FormGroup className={styles.date}>
        <Label htmlFor="date">Date</Label>
        <TextInput className={styles.date} name="date" readOnly value={formatSignatureDate(signatureDate)} />
      </FormGroup>
    </div>
  </div>
);

PrintableLegalese.propTypes = {
  signature: string,
  signatureDate: string,
};

PrintableLegalese.defaultProps = {
  signature: '',
  signatureDate: '',
};

export default PrintableLegalese;
