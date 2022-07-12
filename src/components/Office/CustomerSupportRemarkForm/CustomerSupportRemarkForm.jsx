import React from 'react';
import { connect } from 'react-redux';
import { Button, Textarea, ErrorMessage } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import { Field, Formik } from 'formik';
import { useParams } from 'react-router-dom';
import * as Yup from 'yup';
import classnames from 'classnames';

import customerSupportRemarkFormStyles from './CustomerSupportRemarkForm.module.scss';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import { CUSTOMER_SUPPORT_REMARKS } from 'constants/queryKeys';
import { OfficeUserInfoShape } from 'types/index';
import { selectLoggedInUser } from 'store/entities/selectors';
import { createCustomerSupportRemarkForMove } from 'services/ghcApi';

const CustomerSupportRemarkForm = ({ officeUser }) => {
  const { moveCode } = useParams();

  const [createRemarkMutation] = useMutation(createCustomerSupportRemarkForMove, {
    onSuccess: () => {
      queryCache.invalidateQueries([CUSTOMER_SUPPORT_REMARKS, moveCode]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const initialValues = {
    remark: '',
  };

  const onSubmit = (values, { resetForm, validateForm }) => {
    const body = { content: values.remark, officeUserID: officeUser?.id };
    createRemarkMutation({ locator: moveCode, body });
    resetForm();
    validateForm();
  };

  const validationSchema = Yup.object().shape({
    remark: Yup.string().max(5000, 'Remarks cannot exceed 5000 characters.').required(),
  });

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
      {({ isValid, errors, values }) => {
        const isEmpty = values.remark === '';
        return (
          <Form className={classnames(formStyles.form, customerSupportRemarkFormStyles.remarkForm)}>
            <p className={customerSupportRemarkFormStyles.newRemarkLabel}>
              <small>Use this form to document any customer support provided for this move.</small>
            </p>

            {!isValid && !isEmpty && <ErrorMessage display={!isValid}>{errors.remark}</ErrorMessage>}

            <Field
              as={Textarea}
              label="Add remark"
              name="remark"
              id="remark"
              className={customerSupportRemarkFormStyles.newRemarkTextArea}
              placeholder="Add your remarks here"
              error={!isValid && !isEmpty}
            />

            <Button type="submit" disabled={!isValid}>
              Save
            </Button>
          </Form>
        );
      }}
    </Formik>
  );
};

CustomerSupportRemarkForm.propTypes = {
  officeUser: OfficeUserInfoShape,
};

CustomerSupportRemarkForm.defaultProps = {
  officeUser: {},
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
  };
};

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerSupportRemarkForm);
