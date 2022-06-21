import React, { useRef, useState } from 'react';
import { Button, Textarea, ErrorMessage } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import { Field, Formik } from 'formik';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';
import * as Yup from 'yup';

import customerSupportRemarkStyles from './CustomerSupportRemarkText.module.scss';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { formatCustomerSupportRemarksDate } from 'utils/formatters';
import { CustomerSupportRemarkShape } from 'types/customerSupportRemark';
import Restricted from 'components/Restricted/Restricted';
import { updateCustomerSupportRemarkForMove } from 'services/ghcApi';
import { CUSTOMER_SUPPORT_REMARKS } from 'constants/queryKeys';

const CustomerSupportRemarkText = ({ customerSupportRemark }) => {
  const { moveCode } = useParams();

  const [isCollapsed, setIsCollapsed] = useState(true);
  const [isEdit, setIsEdit] = useState(false);

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

  const toggleEdit = () => {
    setIsEdit(!isEdit);
  };

  const [mutateCustomerSupportRemark] = useMutation(updateCustomerSupportRemarkForMove, {
    onSuccess: async () => {
      await queryCache.invalidateQueries([CUSTOMER_SUPPORT_REMARKS, moveCode]);
      setIsEdit(false);
    },
  });

  const handleSubmitEdit = (values) => {
    mutateCustomerSupportRemark({
      body: {
        id: customerSupportRemark.id,
        content: values.remark,
      },
      locator: moveCode,
    });
  };

  const validationSchema = Yup.object().shape({
    remark: Yup.string().max(5000, 'Remarks cannot exceed 5000 characters.').required(),
  });

  if (isEdit) {
    return (
      <Formik
        initialValues={{ remark: customerSupportRemark.content }}
        onSubmit={handleSubmitEdit}
        validationSchema={validationSchema}
      >
        {({ isValid, errors, values, resetForm }) => {
          const isEmpty = values.remark === '';
          return (
            <div key={customerSupportRemark.id} className={customerSupportRemarkStyles.customerSupportRemarkWrapper}>
              <Form className={classnames(formStyles.form, customerSupportRemarkStyles.remarkForm)}>
                <div className={customerSupportRemarkStyles.row}>
                  <p className={customerSupportRemarkStyles.customerSupportRemarkNameTimestamp}>
                    <small>
                      <strong>
                        {customerSupportRemark.officeUserFirstName} {customerSupportRemark.officeUserLastName}
                      </strong>{' '}
                      {formatCustomerSupportRemarksDate(customerSupportRemark.createdAt)}{' '}
                      {customerSupportRemark.updatedAt !== customerSupportRemark.createdAt && (
                        <small title={`Edited ${formatCustomerSupportRemarksDate(customerSupportRemark.updatedAt)}`}>
                          (edited)
                        </small>
                      )}
                    </small>
                  </p>
                  <Restricted user={customerSupportRemark.officeUserID}>
                    <div className={customerSupportRemarkStyles.row}>
                      <Button
                        className={classnames(
                          customerSupportRemarkStyles.editDeleteButtons,
                          'usa-button',
                          'usa-button--unstyled',
                        )}
                        type="submit"
                        disabled={!isValid}
                        data-testid="edit-remark-save-button"
                      >
                        <small>Save</small>
                      </Button>
                      <small className={customerSupportRemarkStyles.buttonDivider}>|</small>
                      <Button
                        className={classnames(
                          customerSupportRemarkStyles.editDeleteButtons,
                          'usa-button',
                          'usa-button--unstyled',
                        )}
                        type="reset"
                        onClick={() => {
                          toggleEdit();
                          resetForm();
                        }}
                        data-testid="edit-remark-cancel-button"
                      >
                        <small>Cancel</small>
                      </Button>
                    </div>
                  </Restricted>
                </div>

                {!isValid && !isEmpty && <ErrorMessage display={!isValid}>{errors.remark}</ErrorMessage>}

                <Field
                  as={Textarea}
                  label="Edit remark"
                  name="remark"
                  id="remark"
                  className={customerSupportRemarkStyles.editTextArea}
                  placeholder="Add your remarks here"
                  error={!isValid && !isEmpty}
                  data-testid="edit-remark-textarea"
                />
              </Form>
            </div>
          );
        }}
      </Formik>
    );
  }

  return (
    <div key={customerSupportRemark.id} className={customerSupportRemarkStyles.customerSupportRemarkWrapper}>
      <div className={customerSupportRemarkStyles.row}>
        <p className={customerSupportRemarkStyles.customerSupportRemarkNameTimestamp}>
          <small>
            <strong>
              {customerSupportRemark.officeUserFirstName} {customerSupportRemark.officeUserLastName}
            </strong>{' '}
            {formatCustomerSupportRemarksDate(customerSupportRemark.createdAt)}{' '}
            {customerSupportRemark.updatedAt !== customerSupportRemark.createdAt && (
              <small title={`Edited ${formatCustomerSupportRemarksDate(customerSupportRemark.updatedAt)}`}>
                (edited)
              </small>
            )}
          </small>
        </p>
        <Restricted user={customerSupportRemark.officeUserID}>
          <div className={customerSupportRemarkStyles.row}>
            <Button
              className={classnames(
                customerSupportRemarkStyles.editDeleteButtons,
                'usa-button',
                'usa-button--unstyled',
              )}
              type="button"
              onClick={() => {
                toggleEdit();
                setIsCollapsed(true);
              }}
              data-testid="edit-remark-button"
            >
              <small>Edit</small>
            </Button>
            <small className={customerSupportRemarkStyles.buttonDivider}>|</small>
            <Button
              className={classnames(
                customerSupportRemarkStyles.editDeleteButtons,
                'usa-button',
                'usa-button--unstyled',
              )}
              type="delete"
              onClick={() => {}}
              data-testid="delete-remark-button"
            >
              <small>Delete</small>
            </Button>
          </div>
        </Restricted>
      </div>
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
