export namespace models {

	export class Alert {
	    category: string;
	    id: number;
	    name: string;
	    expiryDate: string;
	    daysRemaining: number;
	    severity: string;

	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.expiryDate = source["expiryDate"];
	        this.daysRemaining = source["daysRemaining"];
	        this.severity = source["severity"];
	    }
	}
	export class Antivirus {
	    id: number;
	    name: string;
	    licenseKey: string;
	    computerId?: number;
	    smartphoneId?: number;
	    tabletId?: number;
	    status: string;
	    expiryDate?: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new Antivirus(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.smartphoneId = source["smartphoneId"];
	        this.tabletId = source["tabletId"];
	        this.status = source["status"];
	        this.expiryDate = source["expiryDate"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class AntivirusInput {
	    name: string;
	    licenseKey: string;
	    computerId?: number;
	    smartphoneId?: number;
	    tabletId?: number;
	    status: string;
	    expiryDate?: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new AntivirusInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.smartphoneId = source["smartphoneId"];
	        this.tabletId = source["tabletId"];
	        this.status = source["status"];
	        this.expiryDate = source["expiryDate"];
	        this.notes = source["notes"];
	    }
	}
	export class AppConfig {
	    dbPath: string;
	    density: string;
	    darkMode: boolean;
	    expiryWarningDays: number;

	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dbPath = source["dbPath"];
	        this.density = source["density"];
	        this.darkMode = source["darkMode"];
	        this.expiryWarningDays = source["expiryWarningDays"];
	    }
	}
	export class Computer {
	    id: number;
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new Computer(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class ComputerInput {
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new ComputerInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	    }
	}
	export class DeviceDropdownItem {
	    id: number;
	    name: string;
	    kind: string;

	    static createFrom(source: any = {}) {
	        return new DeviceDropdownItem(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.kind = source["kind"];
	    }
	}
	export class DropdownItem {
	    id: number;
	    name: string;

	    static createFrom(source: any = {}) {
	        return new DropdownItem(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class OtherSoftware {
	    id: number;
	    name: string;
	    licenseKey: string;
	    computerId?: number;
	    smartphoneId?: number;
	    tabletId?: number;
	    status: string;
	    expiryDate?: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new OtherSoftware(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.smartphoneId = source["smartphoneId"];
	        this.tabletId = source["tabletId"];
	        this.status = source["status"];
	        this.expiryDate = source["expiryDate"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class OtherSoftwareInput {
	    name: string;
	    licenseKey: string;
	    computerId?: number;
	    smartphoneId?: number;
	    tabletId?: number;
	    status: string;
	    expiryDate?: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new OtherSoftwareInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.smartphoneId = source["smartphoneId"];
	        this.tabletId = source["tabletId"];
	        this.status = source["status"];
	        this.expiryDate = source["expiryDate"];
	        this.notes = source["notes"];
	    }
	}
	export class Smartphone {
	    id: number;
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new Smartphone(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class SmartphoneInput {
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new SmartphoneInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	    }
	}
	export class Tablet {
	    id: number;
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new Tablet(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class TabletInput {
	    name: string;
	    model: string;
	    userId?: number;
	    status: string;
	    purchaseDate?: string;
	    warrantyExpiry?: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new TabletInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.userId = source["userId"];
	        this.status = source["status"];
	        this.purchaseDate = source["purchaseDate"];
	        this.warrantyExpiry = source["warrantyExpiry"];
	        this.notes = source["notes"];
	    }
	}
	export class User {
	    id: number;
	    name: string;
	    surname?: string;
	    status: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new User(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.surname = source["surname"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class UserInput {
	    name: string;
	    surname?: string;
	    status: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new UserInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.surname = source["surname"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	    }
	}
	export class WindowsKey {
	    id: number;
	    licenseKey: string;
	    computerId?: number;
	    status: string;
	    notes?: string;
	    createdAt: string;
	    updatedAt: string;

	    static createFrom(source: any = {}) {
	        return new WindowsKey(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class WindowsKeyInput {
	    licenseKey: string;
	    computerId?: number;
	    status: string;
	    notes?: string;

	    static createFrom(source: any = {}) {
	        return new WindowsKeyInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.licenseKey = source["licenseKey"];
	        this.computerId = source["computerId"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	    }
	}

}
