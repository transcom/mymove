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
      textBodyRef.current.classList.remove(customerSupportRemarkStyles.customerSupportRemarkTextClamp);
    } else {
      setIsCollapsed(true);
      textBodyRef.current.classList.add(customerSupportRemarkStyles.customerSupportRemarkTextClamp);
    }
  };

  return (
    <div key={customerSupportRemark.id} className={customerSupportRemarkStyles.customerSupportRemarkWrapper}>
      <p className={customerSupportRemarkStyles.customerSupportRemarkNameTimestamp}>
        <small>
          <strong>
            {customerSupportRemark.officeUserFirstName} {customerSupportRemark.officeUserLastName}
          </strong>{' '}
          {formatCustomerSupportRemarksDate(customerSupportRemark.createdAt)}
        </small>
      </p>
      <p
        className={classnames(
          isLongEnough(customerSupportRemark.content) ? customerSupportRemarkStyles.customerSupportRemarkTextClamp : '',
          customerSupportRemarkStyles.customerSupportRemarkText,
        )}
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
