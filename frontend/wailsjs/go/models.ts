export namespace model {
	
	export class HostFileInfo {
	    name: string;
	    isDirty: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HostFileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isDirty = source["isDirty"];
	    }
	}
	export class SaveHostFileRequest {
	    name: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveHostFileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.content = source["content"];
	    }
	}

}

