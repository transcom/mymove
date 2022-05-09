import React, { useRef, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import customerSupportRemarkStyles from './CustomerSupportRemarkText.module.scss';

import { formatCustomerSupportRemarksDate } from 'utils/formatters';
import { CustomerSupportRemarkShape } from 'types/customerSupportRemark';

const CustomerSupportRemarkText = ({ customerSupportRemark }) => {
  const [isCollapsed, setIsCollapsed] = useState(true);

  const isLongEnough = (text) => {
    return text.length >= 350;
  };

  const textBodyRef = useRef(null);

  const handleCollapseClick = () => {
    if (isCollapsed) {
      setIsCollapsed(false);
      textBodyRef.current.classList.remove(customerSupportRemarkStyles.remarksTextContent);
    } else {
      setIsCollapsed(true);
      textBodyRef.current.classList.add(customerSupportRemarkStyles.remarksTextContent);
    }
  };

  return (
    <div key={customerSupportRemark.id}>
      <p className={customerSupportRemarkStyles.customerRemarkBody}>
        <small>
          <strong>
            {customerSupportRemark.officeUserFirstName} {customerSupportRemark.officeUserLastName}
          </strong>{' '}
          {formatCustomerSupportRemarksDate(customerSupportRemark.createdAt)}
        </small>
      </p>
      <p
        className={isLongEnough(customerSupportRemark.content) ? customerSupportRemarkStyles.remarksTextContent : ''}
        ref={textBodyRef}
      >
        {customerSupportRemark.content}
      </p>
      {isLongEnough(customerSupportRemark.content) && (
        <Button
          className={classnames(customerSupportRemarkStyles.seeMoreOrLessButton, 'usa-button', 'usa-button--unstyled')}
          type="button"
          onClick={handleCollapseClick}
        >
          {isCollapsed ? '(see more)' : '(see less)'}
        </Button>
      )}
    </div>
  );
};

CustomerSupportRemarkText.propTypes = {
  customerSupportRemark: CustomerSupportRemarkShape.isRequired,
};

export default CustomerSupportRemarkText;
