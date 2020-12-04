import React, { Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm, SubmissionError } from 'redux-form';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import { validateAccessCode, claimAccessCode } from 'shared/Entities/modules/accessCodes';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

import styles from './AccessCode.module.scss';

const invalidAccessCodeFormatMsg = 'Please check the format';
const invalidAccessCodeMsg = 'This code is invalid';
const claimAccessCodeErrorMsg = 'There was an error. Please reach out to DPS';

class AccessCode extends React.Component {
  validateAccessCodePattern = (code) => {
    const validAccessCodePattern = RegExp('^(HHG|PPM)-[A-Z0-9]{6}$');
    const validAccessCode = validAccessCodePattern.test(code);

    if (!validAccessCode) {
      throw new SubmissionError({
        claim_access_code: invalidAccessCodeFormatMsg,
      });
    }
  };

  validateAndClaimAccessCode = () => {
    const { formValues, validateAccessCode, claimAccessCode } = this.props;
    const code = formValues.claim_access_code;
    this.validateAccessCodePattern(code);

    return validateAccessCode(code)
      .then((res) => {
        const { body: accessCode } = get(res, 'response');
        if (!accessCode.code) {
          throw new SubmissionError({
            claim_access_code: invalidAccessCodeMsg,
          });
        }
        claimAccessCode(accessCode)
          .then(() => {
            window.location.reload();
          })
          .catch((err) => {
            throw new SubmissionError({
              claim_access_code: claimAccessCodeErrorMsg,
            });
          });
      })
      .catch((err) => {
        const errorMsg = get(err, 'errors.claim_access_code');
        throw new SubmissionError({
          claim_access_code: errorMsg,
        });
      });
  };
  render() {
    const { schema, handleSubmit } = this.props;
    return (
      <Fragment>
        <div className="usa-grid">
          <h3 className="title">Welcome to MilMove</h3>
          <p>Please enter your MilMove access code in the field below.</p>
          <SwaggerField
            className={styles['access-code-input']}
            fieldName="claim_access_code"
            swagger={schema}
            required
          />
          <button
            className={classNames('usa-button', styles['submit-access-code'])}
            onClick={handleSubmit(this.validateAndClaimAccessCode)}
          >
            Continue
          </button>
          <br />
          <div className={styles['secondary-text']}>
            No code? Go to{' '}
            <a href="https://eta.sddc.army.mil/ETASSOPortal/default.aspx" className="usa-link">
              DPS
            </a>{' '}
            to schedule your move.
          </div>
        </div>
      </Fragment>
    );
  }
}

const formName = 'claim_access_code_form';
AccessCode = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccessCode);

AccessCode.propTypes = {
  schema: PropTypes.object.isRequired,
  serviceMemberId: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.ClaimAccessCode', {}),
    serviceMemberId: serviceMember?.id,
    formValues: getFormValues(formName)(state),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ validateAccessCode, claimAccessCode }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(AccessCode));
