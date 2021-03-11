import React from 'react';
// import { PropTypes } from 'prop-types';
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ServiceItemCalculations.module.scss';

// import { formatCents, toDollarString } from 'shared/formatters';
// import { PaymentServiceItemShape } from 'types';

const ServiceItemCalculations = () => {
  return (
    <div className={styles.ServiceItemCalculations}>
      <div className={styles.flexGrid}>
        <div className={styles.col}>
          <div className={styles.value}>85 cwt</div>
          <hr />
          <div className={styles.descriptionTitle}>Billable weight (cwt)</div>
          <div className={styles.descriptionContent}>
            Shipment weight: 8,500 lbs <br /> Estimated: 8,000
          </div>
        </div>

        <div className={styles.col}>
          <div className={styles.value}>
            <span className={styles.multiplier}>X</span> 2,337
          </div>
          <hr />
          <div className={styles.descriptionTitle}>Mileage</div>
          <div className={styles.descriptionContent}>Zip 322 to Zip 919</div>
        </div>

        <div className={styles.col}>
          <div className={styles.value}>
            <span className={styles.multiplier}>X</span> 0.03
          </div>
          <hr />
          <div className={styles.descriptionTitle}>Baseline linehaul price</div>
          <div className={styles.descriptionContent}>
            Domestic non-peak <br /> Origin service area: 176 <br /> Pickup date: 24 Jan 2020
          </div>
        </div>

        <div className={styles.col}>
          <div className={styles.value}>
            <span className={styles.multiplier}>X</span> 1.033
          </div>
          <hr />
          <div className={styles.descriptionTitle}>Price escalation factor</div>
          <div className={styles.descriptionContent} />
        </div>

        <div className={styles.col}>
          <div className={styles.value}>
            <span className={styles.equal}>X</span> $6.423
          </div>
          <hr />
          <div className={styles.descriptionTitle}>Total amount request</div>
          <div className={styles.descriptionContent} />
        </div>
      </div>
    </div>
  );
};

// Collection of calculations for the service item
// const ServiceItemCalculations = () => {
//   return (
//     <div className={styles.ServiceItemCalculations}>
//       <div className={styles.flexGrid}>
//         <div className={styles.col}>
//           <div className={styles.value}>85 cwt</div>
//         </div>
//         <div className={styles.col}>
//           <div className={styles.value}>
//             <span className={styles.multiplier}>X</span> 2,337
//           </div>
//         </div>
//         <div className={styles.col}>
//           <div className={styles.value}>
//             <span className={styles.multiplier}>X</span> 0.03
//           </div>
//         </div>
//         <div className={styles.col}>
//           <div className={styles.value}>
//             <span className={styles.multiplier}>X</span> 1.033
//           </div>
//         </div>
//         <div className={styles.col}>
//           <div className={styles.value}>
//             <span className={styles.equal}>X</span> $6.423
//           </div>
//         </div>
//       </div>
//
//       <hr />
//       <div className={styles.flexGrid}>
//         <div className={styles.col}>
//           <div className={styles.descriptionTitle}>Billable weight (cwt)</div>
//           <div className={styles.descriptionContent}>
//             Shipment weight: 8,500 lbs <br /> Estimated: 8,000
//           </div>
//         </div>
//
//         <div className={styles.col}>
//           <div className={styles.descriptionTitle}>Mileage</div>
//           <div className={styles.descriptionContent}>Zip 322 to Zip 919</div>
//         </div>
//
//         <div className={styles.col}>
//           <div className={styles.descriptionTitle}>Baseline linehaul price</div>
//           <div className={styles.descriptionContent}>
//             Domestic non-peak <br /> Origin service area: 176 <br /> Pickup date: 24 Jan 2020
//           </div>
//         </div>
//
//         <div className={styles.col}>
//           <div className={styles.descriptionTitle}>Price escalation factor</div>
//           <div className={styles.descriptionContent} />
//         </div>
//
//         <div className={styles.col}>
//           <div className={styles.descriptionTitle}>Total amount request</div>
//           <div className={styles.descriptionContent} />
//         </div>
//       </div>
//     </div>
//   );
// };

ServiceItemCalculations.propTypes = {};

ServiceItemCalculations.defaultProps = {};

export default ServiceItemCalculations;
