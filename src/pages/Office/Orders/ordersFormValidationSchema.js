import * as Yup from 'yup';

const noStarOrQuote = /^[^*"]*$/;

const ordersFormValidationSchema = Yup.object({
  originDutyLocation: Yup.object().defined('Required'),
  newDutyLocation: Yup.object().required('Required'),
  issueDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  reportByDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  departmentIndicator: Yup.string().required('Required'),
  ordersNumber: Yup.string().required('Required'),
  ordersType: Yup.string().required('Required'),
  ordersTypeDetail: Yup.string().required('Required'),
  tac: Yup.string()
    .matches(noStarOrQuote, 'TAC cannot contain * or " characters')
    .min(4, 'Enter a 4-character TAC')
    .required('Required'),
  sac: Yup.string().matches(noStarOrQuote, 'SAC cannot contain * or " characters'),
  ntsTac: Yup.string().matches(noStarOrQuote, 'NTS TAC cannot contain * or " characters'),
  ntsSac: Yup.string().matches(noStarOrQuote, 'NTS SAC cannot contain * or " characters'),
});

export default ordersFormValidationSchema;
