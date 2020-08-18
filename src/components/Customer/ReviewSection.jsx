// import React from 'react';
// import { Link } from 'react-router-dom';
// // import { get } from 'lodash';
// import PropTypes from 'prop-types';

// // import Address from './Address';

// // import { formatDateSM } from 'shared/formatters';
// // import { getFullSMName } from 'utils/moveSetupFlow';

// import 'scenes/Review/Review.css';

// const serviceMemberDetails = [
//   { label: 'Name', value: 'William Smithingson', key: '123' },
//   { label: 'Rank', value: 'Army', key: '1232' },
//   { label: 'DoD ID#', value: '827309789', key: '12312' },
// ];

// // const sectionTitle = 'Profile';

// const reviewSectionInputs = serviceMemberDetails.map((field) => (
//   <tr key={field.key}>
//     <th scope="row">{field.label}</th>
//     <td>{field.value}</td>
//   </tr>
// ));

// const ReviewSection = ({ fieldData, aTitle, editLink }) => {
//   // console.log('🍐', aTitle);

//   return (
//     <div className="service-member-summary">
//       <div className="stackedtable-header">
//         <div>
//           <h2>
//             {aTitle}
//             <span className="edit-section-link">
//               <Link to={editLink} className="usa-link">
//                 Edit
//               </Link>
//             </span>
//           </h2>
//         </div>
//       </div>
//       <table className="table--stacked review-section">
//         <colgroup>
//           <col style={{ width: '25%' }} />
//           <col style={{ width: '75%' }} />
//         </colgroup>
//         <tbody>{fieldData[0]}</tbody>
//       </table>
//     </div>
//   );
// };

// /*
//   <table className="table--stacked review-section">
//     <colgroup>
//       <col style={{ width: '25%' }} />
//       <col style={{ width: '75%' }} />
//     </colgroup>
//     <tbody>
//       <tr>
//         <th scope="row">Contact info</th>
//       </tr>
//       <tr>
//         <th scope="row">Best contact phone</th>
//         <td>{get(serviceMember, 'telephone')}</td>
//       </tr>
//       <tr>
//         <th scope="row">Personal email</th>
//         <td>{get(serviceMember, 'personal_email')}</td>
//       </tr>
//       <tr>
//         <th scope="row">Current mailing address</th>
//         <td>
//           <Address address={get(serviceMember, 'residential_address')} />
//         </td>
//       </tr>
//     </tbody>
//   </table>
// */

// ReviewSection.propTypes = {
//   fieldData: PropTypes.arrayOf(
//     PropTypes.shape({
//       label: PropTypes.string,
//       value: PropTypes.string,
//       key: PropTypes.string,
//     }),
//   ),
//   aTitle: PropTypes.string,
//   editLink: PropTypes.string,
// };

// export default ReviewSection;
