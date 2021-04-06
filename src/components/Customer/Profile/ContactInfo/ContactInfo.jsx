import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../../Review/Review.module.scss';
import serviceInfoTableStyles from '../../Review/ServiceInfoTable/ServiceInfoTable.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { ResidentialAddressShape } from 'types/address';
import { BackupContactShape } from 'types/customerShapes';

const ContactInfo = ({
  telephone,
  secondaryTelephone,
  personalEmail,
  phoneIsPreferred,
  emailIsPreferred,
  residentialAddress,
  backupMailingAddress,
  backupContact,
  onEditClick,
}) => {
  const containerClassNames = classnames(reviewStyles['review-container'], serviceInfoTableStyles.ServiceInfoTable);
  const tableClassNames = classnames('table--stacked', reviewStyles['review-table']);

  let preferredContactMethod = 'Unknown';
  if (phoneIsPreferred) {
    preferredContactMethod = 'Phone';
  } else if (emailIsPreferred) {
    preferredContactMethod = 'Email';
  }

  return (
    <SectionWrapper className={containerClassNames}>
      <div className={classnames(reviewStyles['review-header'], serviceInfoTableStyles.ReviewHeader)}>
        <h2>Contact info</h2>
        <Button unstyled className={reviewStyles['edit-btn']} data-testid="edit-contact-info" onClick={onEditClick}>
          Edit
        </Button>
      </div>

      <div>
        <table className={tableClassNames}>
          <colgroup>
            <col />
            <col />
          </colgroup>
          <tbody>
            <tr>
              <th scope="row">Best contact phone</th>
              <td>{telephone}</td>
            </tr>
            <tr>
              <th scope="row">Alt. phone</th>
              <td>{secondaryTelephone || '–'}</td>
            </tr>
            <tr>
              <th scope="row">Personal email</th>
              <td>{personalEmail}</td>
            </tr>
            <tr>
              <th scope="row">Preferred contact method</th>
              <td>{preferredContactMethod}</td>
            </tr>
            <tr>
              <th scope="row">Current mailing address</th>
              <td>
                {residentialAddress.street_address_1} {residentialAddress.street_address_2}
                <br />
                {residentialAddress.city}, {residentialAddress.state} {residentialAddress.postal_code}
              </td>
            </tr>
            <tr>
              <th scope="row">Backup mailing address</th>
              <td>
                {backupMailingAddress.street_address_1} {backupMailingAddress.street_address_2}
                <br />
                {backupMailingAddress.city}, {backupMailingAddress.state} {backupMailingAddress.postal_code}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div>
        <h3>Backup contact</h3>
        <table className={tableClassNames}>
          <colgroup>
            <col />
            <col />
          </colgroup>
          <tbody>
            <tr>
              <th scope="row">Name</th>
              <td>{backupContact.name}</td>
            </tr>
            <tr>
              <th scope="row">Email</th>
              <td>{backupContact.email}</td>
            </tr>
            <tr>
              <th className={reviewStyles['table-divider-top']} scope="row" style={{ borderBottom: 'none' }}>
                Phone
              </th>
              <td className={reviewStyles['table-divider-top']} style={{ borderBottom: 'none' }}>
                {backupContact.telephone}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </SectionWrapper>
  );
};

ContactInfo.propTypes = {
  telephone: PropTypes.string.isRequired,
  secondaryTelephone: PropTypes.string,
  personalEmail: PropTypes.string.isRequired,
  phoneIsPreferred: PropTypes.bool.isRequired,
  emailIsPreferred: PropTypes.bool.isRequired,
  residentialAddress: ResidentialAddressShape.isRequired,
  backupMailingAddress: ResidentialAddressShape.isRequired,
  backupContact: BackupContactShape.isRequired,
  onEditClick: PropTypes.func.isRequired,
};

ContactInfo.defaultProps = {
  secondaryTelephone: '',
};

export default ContactInfo;
