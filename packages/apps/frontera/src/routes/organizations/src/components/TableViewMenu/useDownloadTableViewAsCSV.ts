import { useSearchParams } from 'react-router-dom';

import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { ColumnView, ColumnViewType } from '@graphql/types';
import { csvDataMapper as contactCsvDataMapper } from '@organizations/components/Columns/contacts';
import { csvDataMapper as orgCsvDataMapper } from '@organizations/components/Columns/organizations';

enum AdditionalColumnViewType {
  ContactsFirstName = 'CONTACTS_FIRST_NAME',
  ContactsLastName = 'CONTACTS_LAST_NAME',
}

interface IColumnView extends Omit<ColumnView, 'columnType'> {
  columnType: ColumnViewType | AdditionalColumnViewType;
}

const getTableName = (tableViewName: string | undefined) => {
  switch (tableViewName) {
    case 'Targets':
      return 'targets';
    case 'Customers':
      return 'customers';
    case 'Contacts':
      return 'contacts';
    case 'Leads':
      return 'leads';
    case 'Churn':
      return 'churned';
    case 'All orgs':
      return 'organizations';
    default:
      return 'organizations';
  }
};

const convertToCSV = (objArray: Array<Array<string>>): string => {
  return objArray
    .map((row) =>
      row
        .map((cell) => {
          const cleanedCell = `${cell ?? ''}`?.replace(/,/g, '');

          return /[",\n\r]/.test(cleanedCell)
            ? `"${cleanedCell}"`
            : cleanedCell;
        })
        .join(','),
    )
    .join('\r\n');
};

export const useDownloadCsv = () => {
  const store = useStore();
  const [searchParams] = useSearchParams();

  const handleGetData = (): Array<Array<string>> => {
    const preset = searchParams.get('preset');
    const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
    const csvDataMapper =
      tableViewDef?.value.tableType === 'CONTACTS'
        ? contactCsvDataMapper
        : orgCsvDataMapper;

    const visibleColumns = tableViewDef?.value.columns?.filter(
      (column) =>
        column.visible &&
        ![
          ColumnViewType.ContactsAvatar,
          ColumnViewType.OrganizationsAvatar,
        ].includes(column.columnType),
    ) as Array<IColumnView>;

    if (visibleColumns) {
      const nameColumnIndex = visibleColumns.findIndex(
        (column) => column.columnType === ColumnViewType.ContactsName,
      );

      if (nameColumnIndex !== -1) {
        visibleColumns.splice(nameColumnIndex, 1);
        visibleColumns.splice(
          nameColumnIndex,
          0,
          {
            columnId: nameColumnIndex,
            filter: '',
            name: '',
            width: 0,
            columnType: AdditionalColumnViewType.ContactsFirstName,
            visible: true,
          },
          {
            columnType: AdditionalColumnViewType.ContactsLastName,
            visible: true,
            columnId: nameColumnIndex + 1,
            filter: '',
            name: '',
            width: 0,
          },
        );
      }
    }

    const headers = visibleColumns?.map((column) =>
      column.columnType.split('_').join(' '),
    ) as Array<string>;

    const data =
      store.ui.filteredTable?.map((row) => {
        return visibleColumns?.map((column) => {
          const mapper: (d: OrganizationStore | ContactStore) => string =
            csvDataMapper?.[column.columnType as keyof typeof csvDataMapper];
          const rowData = row as ContactStore | OrganizationStore;

          return mapper ? mapper?.(rowData) : '';
        }) as Array<string>;
      }) || [];

    return [headers, ...data] as Array<Array<string>>;
  };

  const downloadCSV = () => {
    const preset = searchParams.get('preset');
    const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
    const tableViewName = tableViewDef?.value.name;
    const tableName = getTableName(tableViewName);

    const data = handleGetData();
    const csvData = new Blob([convertToCSV(data)], { type: 'text/csv' });
    const csvURL = URL.createObjectURL(csvData);
    const link = document.createElement('a');

    link.href = csvURL;
    link.download = `${tableName}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return { downloadCSV };
};
