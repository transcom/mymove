import React from 'react';
import * as PropTypes from 'prop-types';

const AllowancesTable = ({ info }) => {
  return (
    <div>
      <table className="table--stacked">
        <caption>
          <div className="stackedtable-header">
            <h4>Allowances</h4>
          </div>
        </caption>
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row" className="text-bold">
              Branch, rank
            </th>
            <td data-cy="branchRank">{`${info.branch}, ${info.rank}`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Weight allowance
            </th>
            <td data-cy="weightAllowance">{`${info.weightAllowance} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Authorized weight
            </th>
            <td data-cy="authorizedWeight">{`${info.authorizedWeight} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Pro-gear
            </th>
            <td data-cy="progear">{`${info.progear} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Spouse pro-gear
            </th>
            <td data-cy="spouseProgear">{`${info.spouseProgear} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Storage in transit
            </th>
            <td data-cy="storageInTransit">{`${info.storageInTransit} days`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Dependents
            </th>
            <td data-cy="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

AllowancesTable.propTypes = {
  info: PropTypes.shape({
    branch: PropTypes.string,
    rank: PropTypes.string,
    weightAllowance: PropTypes.number,
    authorizedWeight: PropTypes.number,
    progear: PropTypes.number,
    spouseProgear: PropTypes.number,
    storageInTransit: PropTypes.number,
    dependents: PropTypes.bool,
  }).isRequired,
};

export default AllowancesTable;
